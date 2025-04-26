package orm

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/orm/ent"
	"github.com/twiglab/doggy/orm/ent/upload"
	"github.com/twiglab/doggy/pf"
)

type UploadConfig struct {
	MetadataURL string
	Address     string
	Port        int

	User string
	Pwd  string
}

type CameraUpload struct {
	SN     string
	IpAddr string
	Last   time.Time
	UUID1  string
	UUID2  string
}

type EntHandle struct {
	Conf   UploadConfig
	client *ent.Client

	m map[string]pf.CameraSetup

	l sync.Mutex
}

func NewEntHandle(conf UploadConfig, clent *ent.Client) *EntHandle {
	return &EntHandle{
		Conf:   conf,
		client: clent,
		m:      make(map[string]pf.CameraSetup),
	}
}

func (h *EntHandle) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	log.Printf("auto reg sn = %s, ip = %s\n", data.SerialNumber, data.IpAddr)
	device, _ := holo.OpenDevice(data.IpAddr, h.Conf.User, h.Conf.Pwd)
	ids, err := h.deviceId(ctx, device)
	if err != nil {
		return err
	}

	if err = h.metaSub(ctx, device); err != nil {
		return err
	}

	return h.handleUpload(ctx, CameraUpload{
		SN:     data.SerialNumber,
		IpAddr: data.IpAddr,
		Last:   time.Now(),
		UUID1:  ids[0].UUID,
	})
}

func (h *EntHandle) handleUpload(ctx context.Context, u CameraUpload) error {
	err := h.client.Upload.Create().
		SetSn(u.SN).
		SetIP(u.IpAddr).
		SetLastTime(time.Now()).
		SetID1(u.UUID1).
		SetID2(u.UUID2).
		OnConflictColumns(upload.FieldSn).
		UpdateNewValues().
		Exec(ctx)

	if err != nil {
		return err
	}

	h.l.Lock()
	defer h.l.Unlock()

	if _, ok := h.m[u.SN]; !ok {
		h.m[u.SN] = pf.CameraSetup{
			SN:     u.SN,
			IpAddr: u.IpAddr,
			UUID1:  u.UUID1,
			UUID2:  u.UUID2,
		}
	}

	return nil
}

func (h *EntHandle) Reload() error {
	us, err := h.client.Upload.Query().All(context.Background())
	if err != nil {
		return err
	}

	h.l.Lock()
	defer h.l.Unlock()

	clear(h.m)
	for _, u := range us {
		h.m[u.Sn] = pf.CameraSetup{
			SN:     u.Sn,
			IpAddr: u.IP,
			UUID1:  u.ID1,
			UUID2:  u.ID2,
		}
	}

	return nil
}

func (h *EntHandle) deviceId(ctx context.Context, device *holo.Device) ([]holo.DeviceID, error) {
	list, err := device.GetDeviceID(ctx)
	if err != nil {
		return nil, err
	}

	ids := list.IDs
	if len(ids) == 0 {
		return nil, errors.New("no ids")
	}

	return ids, nil
}

func (h *EntHandle) metaSub(ctx context.Context, device *holo.Device) error {
	subscriptions, err := device.GetMetadataSubscription(ctx)
	if err != nil {
		return err
	}

	if !subscriptions.IsEmpty() {
		return nil
	}

	resp, err := device.PostMetadataSubscription(ctx, holo.SubscriptionReq{
		Address:     h.Conf.Address,
		Port:        h.Conf.Port,
		TimeOut:     0,
		HttpsEnable: 1,
		MetadataURL: h.Conf.MetadataURL,
	})

	if err != nil {
		return err
	}

	if resp.IsErr() {
		return resp
	}

	return nil
}
