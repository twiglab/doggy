package pf

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/twiglab/doggy/holo"
)

type AutoSubConf struct {
	MetadataURL string
	Addr        string
	Port        int
}

type CameraUpload struct {
	SN     string
	IpAddr string
	Last   time.Time
	UUID1  string
	UUID2  string
}

type DeviceResolver interface {
	Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error)
}

type FixUserDeviceResolve struct {
	User string
	Pwd  string
}

func (d *FixUserDeviceResolve) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.OpenDevice(data.IpAddr, d.User, d.Pwd)
}

type UploadHandler interface {
	HandleUpload(ctx context.Context, u CameraUpload) error
}

type AutoSub struct {
	Conf           AutoSubConf
	DeviceResolver DeviceResolver
	UploadHandler  UploadHandler
}

func (a *AutoSub) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	log.Printf("auto reg sn = %s, ip = %s\n", data.SerialNumber, data.IpAddr)
	device, err := a.DeviceResolver.Resolve(ctx, data)
	if err != nil {
		return err
	}

	items, err := a.GetSubIDs(ctx, device)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		if err := a.Sub(ctx, device); err != nil {
			return err
		}
	}

	ids, err := a.GetDeviceID(ctx, device)
	if err != nil {
		return err
	}

	return a.UploadHandler.HandleUpload(ctx, CameraUpload{
		SN:     data.SerialNumber,
		IpAddr: data.IpAddr,
		Last:   time.Now(),
		UUID1:  ids[0].UUID,
	})
}

func (a *AutoSub) GetDeviceID(ctx context.Context, device *holo.Device) ([]holo.DeviceID, error) {
	ids, err := device.GetDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	if len(ids.IDs) <= 0 {
		return ids.IDs, errors.New("not found ids")
	}

	return ids.IDs, nil
}

func (a *AutoSub) GetSubIDs(ctx context.Context, device *holo.Device) ([]holo.SubscripionItem, error) {
	subs, err := device.GetMetadataSubscription(ctx)
	if err != nil {
		return nil, err
	}

	return subs.Subscripions, nil
}

func (a *AutoSub) Sub(ctx context.Context, device *holo.Device) error {
	resp, err := device.PostMetadataSubscription(ctx,
		holo.SubscriptionReq{
			Address:     a.Conf.Addr,
			Port:        a.Conf.Port,
			TimeOut:     0,
			HttpsEnable: 1,
			MetadataURL: a.Conf.MetadataURL,
		})
	if err != nil {
		return err
	}

	if resp.IsErr() {
		return resp
	}

	return nil
}
