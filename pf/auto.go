package pf

import (
	"context"
	"time"

	"github.com/twiglab/doggy/holo"
)

type Storer interface {
	Store(ctx context.Context, channels []Channel) error
}

type AutoSub struct {
	DeviceResolver DeviceResolver[*HoloCamera, holo.DeviceAutoRegisterData]
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

	var chs []Channel
	for _, ch := range data.ChannelInfo {
		d, err := camera.ChannelData(ch.UUID)
		if err == nil {
			chs = append(chs, Channel{
				SN:     camera.SerialNumber(),
				IpAddr: camera.IpAddr(),
				UUID:   d.UUID,
				Code:   d.Code,

				X: d.X,
				Y: d.Y,
				Z: d.Z,

				RegTime: time.Now(),
			})
		}
	}
	return a.Storer.Store(ctx, chs)
}
