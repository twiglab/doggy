package pf

import (
	"context"

	"github.com/twiglab/doggy/holo"
)

type CameraDB struct {
	User string
	Pwd  string

	UseHttps bool
}

func (r *CameraDB) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.OpenDevice(data.IpAddr, r.User, r.Pwd, r.UseHttps)
}
