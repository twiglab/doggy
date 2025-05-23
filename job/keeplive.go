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

	MetadataURL string
	Addr        string
	Port        int
}

func (x *KeepLiveJob) Run() {
	ctx := context.Background()
	slog.InfoContext(ctx, "KeepliveJob", slog.String("text", "keepliveJob staring..."))

	ds, err := x.DeviceLoader.All(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "KeepliveJob", slog.String("errText", err.Error()))
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

	l := len(subs.Subscriptions)
	if l != 0 {
		if l > 1 {
			slog.InfoContext(ctx, "KeepliveJob",
				slog.String("sn", data.SN),
				slog.Int("size", l),
				slog.String("errText", "too many subs"),
			)
		}
		return
	}

	resp, err := device.PostMetadataSubscription(ctx, holo.SubscriptionReq{
		Address:     x.Addr,
		Port:        x.Port,
		TimeOut:     0,
		HttpsEnable: 1,
		MetadataURL: x.MetadataURL,
	})

	if err != nil {
		slog.ErrorContext(ctx, "KeepliveJob",
			slog.String("sn", data.SN),
			slog.Int("size", l),
			slog.String("errText", err.Error()),
		)
	}

	if err := resp.Err(); err != nil {
		slog.ErrorContext(ctx, "KeepliveJob",
			slog.String("sn", data.SN),
			slog.Int("size", l),
			slog.String("errText", err.Error()),
		)
	}

	slog.InfoContext(ctx, "KeepliveJob",
		slog.String("sn", data.SN),
		slog.String("ping", "OK"),
	)
}
