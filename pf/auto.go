package pf

import (
	"context"
	"log/slog"
	"time"

	"github.com/twiglab/doggy/holo"
)

type Uploader interface {
	Upload(ctx context.Context, channels []Channel) error
}

type DeviceResolver interface {
	Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error)
}

type AutoSub struct {
	DeviceResolver DeviceResolver
	Uploader       Uploader

	MainSub holo.SubscriptionReq
	Backups []holo.SubscriptionReq
	Muti    int
}

func (a *AutoSub) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {

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

				slog.InfoContext(ctx, "send muti  meta sub",
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

	var chs []Channel
	for _, ch := range data.ChannelInfo {
		chs = append(chs, Channel{
			SN:      data.SerialNumber,
			IpAddr:  data.IpAddr,
			UUID:    ch.UUID,
			Code:    ch.DeviceID,
			RegTime: time.Now(),
		})
	}
	return a.Uploader.Upload(ctx, chs)
}
