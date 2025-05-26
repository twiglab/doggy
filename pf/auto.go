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

	UUID1 string
	Code1 string

	User string
	Pwd  string
}

type DeviceResolver interface {
	Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error)
}

type UploadHandler interface {
	HandleUpload(ctx context.Context, u CameraUpload) error
}

type AutoSub struct {
	DeviceResolver DeviceResolver
	UploadHandler  UploadHandler

	MainSub holo.SubscriptionReq
	Backups []holo.SubscriptionReq
	MutiSub int
}

func (a *AutoSub) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	log.Printf("auto sub reg sn = %s, ip = %s\n", data.SerialNumber, data.IpAddr)
	device, err := a.DeviceResolver.Resolve(ctx, data)
	if err != nil {
		return err
	}
	defer device.Close()

	ids, err := device.GetDeviceID(ctx)
	if err != nil {
		return err
	}

	if len(ids.IDs) < 1 {
		return errors.New("not found device ids")
	}

	id := ids.IDs[0]
	if device.DeviceID != "" {
		if device.DeviceID != id.DeviceID {

			log.Println("下发DeviceID ", device.DeviceID)

			res, err := device.PutDeviceID(ctx,
				holo.DeviceIDList{
					IDs: []holo.DeviceID{
						{UUID: id.UUID, DeviceID: device.DeviceID},
					},
				})

			if err != nil {
				return err
			}
			if err := res.Err(); err != nil {
				return err
			}
		}
	} else {
		device.DeviceID = id.DeviceID
		device.UUID = id.UUID
	}

	subs, err := device.GetMetadataSubscription(ctx)
	if err != nil {
		return err
	}

	if len(subs.Subscriptions) == 0 {

		log.Println("下发元数据订阅参数", a.MainSub.MetadataURL)

		res, err := device.PostMetadataSubscription(ctx, a.MainSub)
		if err := holo.CheckErr(res, err); err != nil {
			return err
		}

		if a.MutiSub != 0 {
			for _, sub := range a.Backups {
				res, err := device.PostMetadataSubscription(ctx, sub)
				if err := holo.CheckErr(res, err); err != nil {
					return err
				}
			}
		}
	}

	return a.UploadHandler.HandleUpload(ctx, CameraUpload{
		SN:     data.SerialNumber,
		IpAddr: data.IpAddr,
		Last:   time.Now(),
		UUID1:  device.UUID,
		Code1:  device.DeviceID,
		User:   device.User,
		Pwd:    device.Pwd,
	})
}
