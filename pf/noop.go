package pf

import (
	"context"
	"log/slog"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/pkg/human"
)

type noopAction struct {
}

func (d *noopAction) Name() string {
	return "noop"
}

func (d noopAction) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	slog.DebugContext(ctx, "receive reg data",
		slog.String("module", "cameraAction"),
		slog.String("sn", data.SerialNumber),
		slog.String("addr", data.IpAddr))
	return nil
}

func (d noopAction) HandleData(ctx context.Context, data human.DataMix) error {
	slog.DebugContext(ctx, "handleData", slog.Any("data", data))
	return nil
}

var NoopAction = noopAction{}
