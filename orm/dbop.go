package orm

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/orm/ent"
)

type Config struct {
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
	Conf   Config
	client *ent.Client
}

func NewEntHandle(clent *ent.Client) *EntHandle {
	return &EntHandle{client: clent}
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

func (h *EntHandle) handleUpload(cx context.Context, upload CameraUpload) error {
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
