package pf

import (
	"context"
	"log"

	"github.com/twiglab/doggy/holo"
)

type KeepLiveJob struct {
	DeviceLoader   DeviceLoader
	DeviceResolver DeviceResolver

	MetadataURL string
	Addr        string
	Port        int
}

func (x *KeepLiveJob) Run() {
	ctx := context.Background()

	ds, err := x.DeviceLoader.All(ctx)
	if err != nil {
		return
	}

	for _, d := range ds {
		go x.Ping(ctx, d)
	}
}

func (x *KeepLiveJob) Ping(ctx context.Context, data CameraUpload) {
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
			log.Printf("%s %s %d too many subs", data.SN, data.IpAddr, l)
		}
		return
	}

	device.PostMetadataSubscription(ctx, holo.SubscriptionReq{
		Address:     x.Addr,
		Port:        x.Port,
		TimeOut:     0,
		HttpsEnable: 1,
		MetadataURL: x.MetadataURL,
	})
}
