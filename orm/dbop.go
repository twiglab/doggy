package orm

import (
	"context"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/orm/ent"
	"github.com/twiglab/doggy/orm/ent/autoreg"
)

type EntHandle struct {
	client *ent.Client
}

func NewEntHandle(clent *ent.Client) *EntHandle {
	return &EntHandle{client: clent}
}

func (h *EntHandle) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	return h.client.AutoReg.Create().
		SetSn(data.SerialNumber).
		SetIP(data.IpAddr).OnConflictColumns(autoreg.FieldSn).
		UpdateNewValues().Exec(ctx)
}
