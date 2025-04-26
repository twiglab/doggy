package orm

import (
	"context"
	"maps"
	"sync"

	"github.com/twiglab/doggy/orm/ent"
	"github.com/twiglab/doggy/pf"
)

type EntUsing struct {
	client *ent.Client

	m map[string]pf.CameraUsing

	mu sync.Mutex
}

func (e *EntUsing) Load() error {
	us, err := e.client.Using.Query().All(context.Background())
	if err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	clear(e.m)
	var bk string
	for _, u := range us {
		bk = NewBK(u.UUID, u.Alg)
		e.m[bk] = pf.CameraUsing{
			SN:    u.Sn,
			UUID:  u.UUID,
			AlgID: u.Alg,
			Name:  u.Name,
			BK:    bk,
		}
	}
	return nil
}

func (e *EntUsing) GetCameraUsings(context.Context) map[string]pf.CameraUsing {
	e.mu.Lock()
	defer e.mu.Unlock()

	return maps.Clone(e.m)

}

func NewBK(uuid, alg string) string {
	return alg + "_" + uuid
}
