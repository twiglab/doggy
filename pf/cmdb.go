package pf

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/twiglab/doggy/holo"
)

type CameraDB struct {
	User string
	Pwd  string

	UseHttps bool
}

func (r *CameraDB) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.OpenDevice(data.IpAddr, r.User, r.Pwd, true)
}

type Toucher interface {
	Last(me string) (t time.Time, ok bool)
	Touch(me string) error
}

type InMomoryTouch struct {
	mm       sync.Map
	zeroTime time.Time
}

func (p *InMomoryTouch) Last(me string) (time.Time, bool) {
	return p.zeroTime, false
}

func (p *InMomoryTouch) Touch(me string) error {
	if me != "" {
		p.mm.Store(me, time.Now())
	}
	return nil
}

type Cache interface {
	Get(context.Context, string) (Channel, bool, error)
	Set(context.Context, Channel) error
}

type emptyCache string

func (i emptyCache) Get(_ context.Context, _ string) (c Channel, ok bool, err error) {
	return
}

func (i emptyCache) Set(_ context.Context, _ Channel) (err error) {
	return
}

type TiersCache struct {
	m      sync.Map
	second Cache
}

func NewTiersCache() *TiersCache {
	return &TiersCache{
		second: emptyCache("x"),
	}
}

func (c *TiersCache) SetSecond(second Cache) {
	c.second = second
}

func (c *TiersCache) Get(ctx context.Context, k string) (i Channel, ok bool, err error) {
	if i, ok, err = c.innerGet(k); ok {
		return
	}

	if i, ok, err = c.second.Get(ctx, k); ok {
		err = c.innerSet(i)
		slog.DebugContext(ctx, "second cache",
			slog.Any("camera", i),
			slog.Any("setCacheErr", err),
			slog.Bool("ok", ok))
	}

	return
}

func (c *TiersCache) Set(ctx context.Context, item Channel) error {
	return c.innerSet(item)
}

func (c *TiersCache) innerGet(k string) (Channel, bool, error) {
	a, ok := c.m.Load(k)
	if ok {
		return a.(Channel), ok, nil
	}

	return Channel{}, false, nil
}

func (c *TiersCache) innerSet(item Channel) error {
	c.m.Store(item.UUID, item)
	return nil
}
