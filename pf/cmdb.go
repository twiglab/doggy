package pf

import (
	"context"
	"sync"
	"time"

	"github.com/twiglab/doggy/holo"
)

type CameraUpload struct {
	SN       string
	IpAddr   string
	LastTime time.Time

	UUID string
	Code string

	User string
	Pwd  string
}

type CameraDB struct {
	User string
	Pwd  string
}

func (r *CameraDB) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.ConnectDevice(data.IpAddr, r.User, r.Pwd)
}

type Toucher interface {
	Last(me string) (t time.Time, ok bool)
	Touch(me string) error
}

type InMomoryTouch struct {
	mm       sync.Map
	zeroTime time.Time
}

func NewInMomoryTouch() *InMomoryTouch {
	return &InMomoryTouch{
		zeroTime: time.Unix(0, 0),
	}
}

func (p *InMomoryTouch) Last(me string) (time.Time, bool) {
	if me != "" {
		if v, ok := p.mm.Load(me); ok {
			return v.(time.Time), ok
		}
	}
	return p.zeroTime, false
}

func (p *InMomoryTouch) Touch(me string) error {
	if me != "" {
		p.mm.Store(me, time.Now())
	}
	return nil
}
