package pf

import (
	"context"
	"log/slog"
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

/*
func (a *AutoSub) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	slog.InfoContext(ctx, "receive reg data",
		slog.String("module", "AutoSub"),
		slog.String("sn", data.SerialNumber),
		slog.String("addr", data.IpAddr))

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

			slog.InfoContext(ctx, "send device id",
				slog.String("deviceID", device.DeviceID),
				slog.String("sn", data.SerialNumber),
				slog.String("module", "AutoSub"),
				slog.String("method", "AutoRegister"))

			res, err := device.PutDeviceID(ctx,
				holo.DeviceIDList{
					IDs: []holo.DeviceID{
						{UUID: id.UUID, DeviceID: device.DeviceID},
					},
				})

			if err := holo.CheckErr(res, err); err != nil {
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

		slog.InfoContext(ctx, "send meta sub",
			slog.Any("main", a.MainSub),
			slog.String("sn", data.SerialNumber),
			slog.String("module", "AutoSub"),
			slog.String("method", "AutoRegister"))

		res, err := device.PostMetadataSubscription(ctx, a.MainSub)
		if err := holo.CheckErr(res, err); err != nil {
			return err
		}

		if a.MutiSub != 0 {
			slog.InfoContext(ctx, "send meta sub",
				slog.Any("backups", a.Backups),
				slog.String("sn", data.SerialNumber),
				slog.String("module", "AutoSub"),
				slog.String("method", "AutoRegister"))
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
*/

func (a *AutoSub) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	/*
		slog.InfoContext(ctx, "receive reg data",
			slog.String("module", "AutoSub"),
			slog.String("sn", data.SerialNumber),
			slog.String("addr", data.IpAddr))
	*/

	go a.reg(ctx, data)

	return nil
}

func (a *AutoSub) reg(ctx context.Context, data holo.DeviceAutoRegisterData) {
	slog.InfoContext(ctx, "receive reg data",
		slog.String("module", "AutoSub"),
		slog.String("sn", data.SerialNumber),
		slog.String("addr", data.IpAddr))

	device, err := a.DeviceResolver.Resolve(ctx, data)
	if err != nil {
		slog.ErrorContext(ctx, "device resolver error",
			slog.String("sn", data.SerialNumber),
			slog.Any("error", err),
		)
		return
	}
	defer device.Close()

	ids, err := device.GetDeviceID(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "GetDeviceID error",
			slog.String("sn", data.SerialNumber),
			slog.Any("error", err),
		)
		return
	}

	if len(ids.IDs) < 1 {
		slog.ErrorContext(ctx, "not found device ids",
			slog.String("sn", data.SerialNumber),
		)
		return
	}

	id := ids.IDs[0]
	if device.DeviceID != "" {
		if device.DeviceID != id.DeviceID {

			slog.InfoContext(ctx, "send device id",
				slog.String("deviceID", device.DeviceID),
				slog.String("sn", data.SerialNumber),
				slog.String("module", "AutoSub"),
				slog.String("method", "AutoRegister"))

			res, err := device.PutDeviceID(ctx,
				holo.DeviceIDList{
					IDs: []holo.DeviceID{
						{UUID: id.UUID, DeviceID: device.DeviceID},
					},
				})

			if err := holo.CheckErr(res, err); err != nil {
				slog.ErrorContext(ctx, "PutDeviceID error",
					slog.String("sn", data.SerialNumber),
					slog.Any("error", err),
				)
				return
			}
		}
	} else {
		device.DeviceID = id.DeviceID
		device.UUID = id.UUID
	}

	subs, err := device.GetMetadataSubscription(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "GetMetadata error",
			slog.String("sn", data.SerialNumber),
			slog.Any("error", err),
		)
		return
	}

	if len(subs.Subscriptions) == 0 {

		slog.InfoContext(ctx, "send meta main sub",
			slog.Any("main", a.MainSub),
			slog.String("sn", data.SerialNumber),
			slog.String("module", "AutoSub"),
			slog.String("method", "AutoRegister"))

		res, err := device.PostMetadataSubscription(ctx, a.MainSub)
		if err := holo.CheckErr(res, err); err != nil {
			slog.ErrorContext(ctx, "PostMetadata main error",
				slog.String("sn", data.SerialNumber),
				slog.Any("meta", a.MainSub),
				slog.Any("error", err),
			)
			return
		}

		if a.MutiSub != 0 {
			slog.InfoContext(ctx, "send meta backups sub",
				slog.Any("backups", a.Backups),
				slog.String("sn", data.SerialNumber),
				slog.String("module", "AutoSub"),
				slog.String("method", "AutoRegister"))
			for _, sub := range a.Backups {
				res, err := device.PostMetadataSubscription(ctx, sub)
				if err := holo.CheckErr(res, err); err != nil {
					slog.ErrorContext(ctx, "PostMetadata backup error",
						slog.String("sn", data.SerialNumber),
						slog.Any("meta", sub),
						slog.Any("error", err),
					)
					return
				}
			}
		}
	}

	err = a.UploadHandler.HandleUpload(ctx, CameraUpload{
		SN:     data.SerialNumber,
		IpAddr: data.IpAddr,
		Last:   time.Now(),
		UUID1:  device.UUID,
		Code1:  device.DeviceID,
		User:   device.User,
		Pwd:    device.Pwd,
	})

	if err != nil {
		slog.ErrorContext(ctx, "Upload error",
			slog.String("sn", data.SerialNumber),
			slog.Any("error", err),
		)
	}
}
