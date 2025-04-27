package pf

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/twiglab/doggy/holo"
)

type CameraUpload struct {
	SN     string
	IpAddr string
	Last   time.Time
	UUID1  string
	UUID2  string

	User string
	Pwd  string
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
	DeviceResolver DeviceResolver
	UploadHandler  UploadHandler

	MetadataURL string
	Addr        string
	Port        int
}

func (a *AutoSub) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	log.Printf("auto reg sn = %s, ip = %s\n", data.SerialNumber, data.IpAddr)
	device, err := a.DeviceResolver.Resolve(ctx, data)
	if err != nil {
		return err
	}
	defer device.Close()

	subs, err := device.GetMetadataSubscription(ctx)
	if err != nil {
		return err
	}

	if len(subs.Subscripions) == 0 {
		_, err := device.PostMetadataSubscription(ctx,
			holo.SubscriptionReq{
				Address:     a.Addr,
				Port:        a.Port,
				TimeOut:     0,
				HttpsEnable: 1,
				MetadataURL: a.MetadataURL,
			})
		if err != nil {
			return err
		}
	}

	ids, err := device.GetDeviceID(ctx)
	if err != nil {
		return err
	}

	if len(ids.IDs) <= 0 {
		return errors.New("not found device ids")
	}

	return a.UploadHandler.HandleUpload(ctx, CameraUpload{
		SN:     data.SerialNumber,
		IpAddr: data.IpAddr,
		Last:   time.Now(),
		UUID1:  ids.IDs[0].UUID,
		User:   device.User,
		Pwd:    device.Pwd,
	})
}
