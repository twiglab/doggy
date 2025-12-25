package pf

import (
	"context"
	"log/slog"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/pkg/human"
)

type LogAction struct {
	log *slog.Logger
}

func NewLogAction(log *slog.Logger) LogAction {
	l := slog.Default()
	if log != nil {
		l = log
	}
	return LogAction{log: l}
}

func (d LogAction) Name() string {
	return "log"
}

func (d LogAction) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	slog.DebugContext(ctx, "receive reg data",
		slog.String("module", "cameraAction"),
		slog.String("sn", data.SerialNumber),
		slog.String("addr", data.IpAddr))
	return nil
}

func (d LogAction) HandleData(ctx context.Context, data human.DataMix) error {
	slog.DebugContext(ctx, "handleData", slog.Any("data", data))
	return nil
}
