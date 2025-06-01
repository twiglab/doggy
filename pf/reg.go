package pf

import (
	"context"
	"log/slog"
	"time"

	"github.com/twiglab/doggy/holo"
)

type AutoReg struct {
	DeviceResolver DeviceResolver
	UploadHandler  UploadHandler
}

func (a *AutoReg) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	device, err := a.DeviceResolver.Resolve(ctx, data)
	if err != nil {
		slog.ErrorContext(ctx, "device resolver error",
			slog.String("sn", data.SerialNumber),
			slog.Any("error", err),
		)
		return err
	}
	defer device.Close()

	return a.UploadHandler.HandleUpload(ctx, CameraUpload{
		SN:     data.SerialNumber,
		IpAddr: data.IpAddr,
		Last:   time.Now(),
		User:   device.User,
		Pwd:    device.Pwd,
	})
}
