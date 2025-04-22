package orm

import (
	"context"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/orm/ent"
)

type EntHandle struct {
	client *ent.Client
}

func NewEntHandle(clent *ent.Client) *EntHandle {
	return &EntHandle{client: clent}
}

func (h *EntHandle) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	return nil
}
