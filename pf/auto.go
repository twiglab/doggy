package pf

import (
	"context"
	"log/slog"
	"time"

	"github.com/twiglab/doggy/holo"
)

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
	Muti    int
}

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

	if a.Muti > 0 {
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

			if a.Muti > 1 {

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
		} else {
			slog.InfoContext(ctx, "metadata size > 0",
				slog.String("sn", data.SerialNumber),
				slog.String("module", "AutoSub"),
				slog.String("method", "AutoRegister"))
		}
	}

	ch := data.FirstChannel()
	return a.UploadHandler.HandleUpload(ctx, CameraUpload{
		SN:       data.SerialNumber,
		IpAddr:   data.IpAddr,
		UUID:     ch.UUID,
		Code:     ch.DeviceID,
		LastTime: time.Now(),
		User:     device.User,
		Pwd:      device.Pwd,
	})
}
