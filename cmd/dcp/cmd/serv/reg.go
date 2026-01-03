package serv

import (
	"context"

	"github.com/spf13/viper"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/kv"
	"github.com/twiglab/doggy/pf"
)

func pfHandle2(ctx context.Context, v *viper.Viper) (pf.DeviceRegister, context.Context) {
	kvh := ctx.Value(keyKvHandle).(*kv.Handle)
	backend := ctx.Value(keyBackend).(pf.DataHandler)

	var backups []holo.SubscriptionReq
	bs := v.GetStringSlice("camera.setup.backups")
	for _, b := range bs {
		backups = append(backups, MustSubReq(holo.SubReq(b)))
	}

	autoSub := &pf.AutoSub{
		DeviceResolver: cmdb,
		Storer:         &kv.Store{H: kvh},
		/*
			MainSub:        MustSubReq(holo.SubReq(conf.SubsConf.Main)),
			Backups:        backups,
			Muti:           conf.SubsConf.Muti,
		*/
	}

	cache := pf.NewTiersCache[string, pf.ChannelExtra]().WithSecond(&kv.ChannelCache{H: kvh})
	toucher := kv.NewTouch(kvh, 90)

	h := pf.NewMainHandle(
		conf.ProjectConf.Project,

		pf.WithDeviceRegister(autoSub),
		pf.WithDataHandler(backend),
		pf.WithToucher(toucher),
		pf.WithCache(cache),
	)

	return pf.PlatformHandle(h)
}
