package db

import (
	"context"

	"github.com/twiglab/doggy/db/ent"
	"github.com/twiglab/doggy/db/ent/autoreg"
	"github.com/twiglab/doggy/holo"
)

type DBOP struct {
	client *ent.Client
}

func (op *DBOP) Reg(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	return op.client.AutoReg.Create().
		SetSn(data.SerialNumber).
		SetIP(data.IpAddr).OnConflictColumns(autoreg.FieldSn).
		Exec(ctx)
}
