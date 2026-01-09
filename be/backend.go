package be

import (
	"context"
	"errors"
	"log/slog"

	"github.com/twiglab/doggy/pkg/human"
)

var ErrUnimplType = errors.New("unsupport type")

const (
	TAOS = "taos"
	MQTT = "mqtt" // mqtt 3.11
	LOG  = "log"

	MQTT5 = "mqtt5" // mqtt5 保留
	HTTP  = "http"  // 保留，暂不提供
	NOOP  = "noop"  // 保留，仅作占位符
)

func HasHuman(in, out int) bool {
	return in != 0 || out != 0
}

func HasCount(count int) bool {
	return count != 0
}

type DataHandler interface {
	HandleData(ctx context.Context, data human.DataMix) error
	Name() string
}

type MutiAction struct {
	actions []DataHandler
}

func NewMutiAction(actions ...DataHandler) *MutiAction {
	return &MutiAction{actions: actions}
}

func (a MutiAction) Name() string {
	return "muti"
}

func (a *MutiAction) Add(h DataHandler) *MutiAction {
	a.actions = append(a.actions, h)
	return a
}

func (a MutiAction) HandleData(ctx context.Context, data human.DataMix) error {
	for _, action := range a.actions {
		if err := action.HandleData(ctx, data); err != nil {
			return err
		}
	}
	return nil
}

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
	return LOG
}

func (d LogAction) HandleData(ctx context.Context, data human.DataMix) error {
	return nil
}
