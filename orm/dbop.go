package orm

import (
	"context"

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

func (h *EntHandle) HandleUpload(ctx context.Context, u pf.CameraItem) error {
	err := h.client.Upload.Create().
		SetSn(u.SN).
		SetIPAddr(u.IpAddr).
		SetUUID(u.UUID).
		SetCode(u.Code).
		SetRegTime(u.RegTime).
		OnConflictColumns(upload.FieldSn).
		UpdateNewValues().
		Exec(ctx)
	return err
}

func (h *EntHandle) All(ctx context.Context) ([]*ent.Upload, error) {
	return h.client.Upload.Query().All(ctx)
}

func (h *EntHandle) GetByUUID(ctx context.Context, uuid string) (*ent.Upload, error) {
	q := h.client.Upload.Query()
	q.Where(upload.UUIDEQ(uuid))
	return q.Only(ctx)
}

func (h *EntHandle) Get(ctx context.Context, k string) (pf.CameraItem, bool, error) {
	u, err := h.GetByUUID(ctx, k)
	if err != nil {
		return pf.CameraItem{}, false, err
	}

	return pf.CameraItem{
		SN:      u.Sn,
		IpAddr:  u.IPAddr,
		UUID:    u.UUID,
		Code:    u.Code,
		RegTime: u.RegTime,
	}, true, nil
}

func (h *EntHandle) Set(ctx context.Context, u pf.CameraItem) error {
	return h.HandleUpload(ctx, u)
}

func (h *EntHandle) GetAll(ctx context.Context) (uploads []pf.CameraItem, err error) {
	us, err := h.All(ctx)
	if err != nil {
		return
	}

	for _, u := range us {
		uploads = append(uploads, pf.CameraItem{
			SN:      u.Sn,
			IpAddr:  u.IPAddr,
			UUID:    u.UUID,
			Code:    u.Code,
			RegTime: u.RegTime,
		})
	}
	return
}
