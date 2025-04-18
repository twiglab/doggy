package orm

import (
	"context"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/orm/ent"
	"github.com/twiglab/doggy/orm/ent/autoreg"
)

type DBOP struct {
	client *ent.Client
}

func NewDBOP(clent *ent.Client) *DBOP {
	return &DBOP{client: clent}
}

func (op *DBOP) Reg(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	return op.client.AutoReg.Create().
		SetSn(data.SerialNumber).
		SetIP(data.IpAddr).OnConflictColumns(autoreg.FieldSn).
		UpdateNewValues().Exec(ctx)
}
