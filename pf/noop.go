package pf

import (
	"context"
	"log/slog"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/pkg/human"
)

type NoopAction struct{}

func (d *NoopAction) Name() string {
	return "noop"
}

func (d *NoopAction) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	slog.DebugContext(ctx, "receive reg data",
		slog.String("module", "cameraAction"),
		slog.String("sn", data.SerialNumber),
		slog.String("addr", data.IpAddr))
	return nil
}

func (d *NoopAction) HandleData(ctx context.Context, data human.DataMix) error {
	slog.DebugContext(ctx, "handleData", slog.Any("data", data))
	return nil
}
