package orm

import (
	"context"
	"time"

	"github.com/twiglab/doggy/orm/ent"
	"github.com/twiglab/doggy/orm/ent/upload"
	"github.com/twiglab/doggy/pf"
)

type EntHandle struct {
	client *ent.Client
}

func NewEntHandle(client *ent.Client) *EntHandle {
	return &EntHandle{
		client: client,
	}
}

func (h *EntHandle) HandleUpload(ctx context.Context, u pf.CameraUpload) error {
	err := h.client.Upload.Create().
		SetSn(u.SN).
		SetIP(u.IpAddr).
		SetLastTime(time.Now()).
		SetID1(u.UUID1).
		SetID2(u.UUID2).
		SetUser(u.User).
		SetPwd(u.Pwd).
		OnConflictColumns(upload.FieldSn).
		UpdateNewValues().
		Exec(ctx)
	return err
}

func (h *EntHandle) All(ctx context.Context) ([]pf.CameraUpload, error) {
	us, err := h.client.Upload.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	uploads := make([]pf.CameraUpload, len(us))
	for i, u := range us {
		uploads[i] = pf.CameraUpload{
			SN:     u.Sn,
			IpAddr: u.IP,
			User:   u.User,
			Pwd:    u.Pwd,
		}
	}
	return uploads, nil
}
