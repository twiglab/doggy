package be

import (
	"context"

	"github.com/twiglab/doggy/pkg/human"
)

const (
	TAOS  = "taos"
	MQTT  = "mqtt"
	MQTT5 = "mqtt5"
	NONE  = "none"
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

func (a *MutiAction) HandleData(ctx context.Context, data human.DataMix) error {
	for _, action := range a.actions {
		if err := action.HandleData(ctx, data); err != nil {
			return err
		}
	}
	return nil
}
