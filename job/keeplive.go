package job

import (
	"context"
	"log/slog"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/pf"
)

type KeepLiveJob struct {
	DeviceLoader   pf.DeviceLoader
	DeviceResolver pf.DeviceResolver

	MainSub holo.SubscriptionReq
}

func (x *KeepLiveJob) Run() {
	ctx := context.Background()
	slog.InfoContext(ctx, "KeepliveJob starting...")

	ds, err := x.DeviceLoader.All(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "KeepliveJob", slog.Any("error", err))
		return
	}

	for _, d := range ds {
		go x.Ping(ctx, d)
	}
}

func (x *KeepLiveJob) Ping(ctx context.Context, data pf.CameraUpload) {
	device, err := x.DeviceResolver.Resolve(ctx,
		holo.DeviceAutoRegisterData{
			SerialNumber: data.SN,
			IpAddr:       data.IpAddr,
		})
	if err != nil {
		return
	}

	subs, err := device.GetMetadataSubscription(ctx)
	if err != nil {
		return
	}

	size := len(subs.Subscriptions)
	if size != 0 {
		slog.InfoContext(ctx, "KeepliveJob ok",
			slog.String("sn", data.SN),
			slog.Int("size", size),
		)
		return
	}

	resp, err := device.PostMetadataSubscription(ctx, x.MainSub)

	if err := holo.CheckErr(resp, err); err != nil {
		slog.ErrorContext(ctx, "KeepliveJob error",
			slog.String("sn", data.SN),
			slog.Int("size", size),
			slog.Any("error", err),
		)
		return
	}

	slog.InfoContext(ctx, "KeepliveJob no sub",
		slog.String("sn", data.SN),
	)
}
