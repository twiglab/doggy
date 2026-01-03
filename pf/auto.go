package pf

import (
	"context"
	"time"

	"github.com/twiglab/doggy/holo"
)

type Storer interface {
	Store(ctx context.Context, channels []ChannelExtra) error
}

type AutoSub struct {
	DeviceResolver DeviceResolver[holo.DeviceAutoRegisterData, *HoloCamera]
	Storer         Storer
}

func (a *AutoSub) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {

	camera, err := a.DeviceResolver.Resolve(ctx, data)
	if err != nil {
		return err
	}
	defer camera.Close()

	if err := camera.Setup(ctx); err != nil {
		return err
	}

	var chs []ChannelExtra
	for _, ch := range data.ChannelInfo {
		if d, err := camera.ChannelData(ctx, ch.UUID); err == nil { // 无错继续， 有错跳过
			chs = append(chs, ChannelExtra{
				SN:     camera.SerialNumber(),
				IpAddr: camera.IpAddr(),

				UUID: d.UUID,
				Code: d.Code,

				X: d.X,
				Y: d.Y,
				Z: d.Z,

				RegTime: time.Now(),
			})
		}
	}
	return a.Storer.Store(ctx, chs)
}
