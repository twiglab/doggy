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
		chs = append(chs, Channel{
			SN:      camera.SerialNumber(),
			IpAddr:  camera.IpAddr(),
			UUID:    ch.UUID,
			Code:    ch.DeviceID,
			RegTime: time.Now(),
		})
	}
	return a.Storer.Store(ctx, chs)
}
