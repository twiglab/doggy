package be

import (
	"context"
	"errors"

	"github.com/twiglab/doggy/pkg/human"
)

var ErrUnimplType = errors.New("unsupport type")

const (
	TAOS  = "taos"
	MQTT  = "mqtt"  // mqtt 3.11
	MQTT5 = "mqtt5" // mqtt5 保留
	NOOP  = "noop"
	LOG   = "log"
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
