package serv

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"github.com/twiglab/doggy/ddb"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/kv"
	"github.com/twiglab/doggy/pf"
)

func buildReg(ctx context.Context, v *viper.Viper) (pf.DeviceRegister, context.Context) {
	kvh := ctx.Value(keyKvHandle).(*kv.Handle)
	//backend := ctx.Value(keyBackend).(pf.DataHandler)

	db, err := ddb.New(v.GetString("cmdb.ddb.from"))
	if err != nil {
		log.Fatal(err)
	}

	main := MustSubReq(holo.SubReq(v.GetString("camera.setup.main")))
	/*
		var backups []holo.SubscriptionReq
		bs := v.GetStringSlice("camera.setup.backups")
		for _, b := range bs {
			backups = append(backups, MustSubReq(holo.SubReq(b)))
		}
	*/

	cameraDB := &pf.CameraDB{
		User:     v.GetString("camera.user"),
		Pwd:      v.GetString("camera.pwd"),
		UseHttps: true,

		Setup: pf.HoloCameraSetup{
			Muti:    1,
			MainSub: main,
			// Backups: backups,
		},

		UserData: db,
	}

	autoSub := &pf.AutoSub{
		DeviceResolver: cameraDB,
		Storer:         &kv.Store{H: kvh},
	}

	//cache := pf.NewTiersCache[string, pf.ChannelExtra]().WithSecond(&kv.ChannelCache{H: kvh})
	//toucher := kv.NewTouch(kvh, 90)

	return autoSub, context.WithValue(ctx, keyReg, autoSub)
	/*

		h := pf.NewMainHandle(
			v.GetString("project"),

			pf.WithDeviceRegister(autoSub),
			pf.WithDataHandler(backend),
			pf.WithToucher(toucher),
			pf.WithCache(cache),
		)

		return pf.PlatformHandle(h)
	*/
}
