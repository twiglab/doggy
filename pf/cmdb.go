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

type CacheSetter interface {
	Set(context.Context, CameraItem) error
}

type CacheGetter interface {
	Get(context.Context, string) (CameraItem, bool, error)
}

type Cache interface {
	CacheGetter
	CacheSetter
}

type CameraItem struct {
	SN       string
	IpAddr   string
	LastTime time.Time

	UUID string
	Code string

	User string
	Pwd  string
}

type emptyCache string

func (i emptyCache) Get(_ context.Context, _ string) (c CameraItem, ok bool, err error) {
	return
}

func (i emptyCache) Set(_ context.Context, _ CameraItem) (err error) {
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

func (c *TiersCache) Get(ctx context.Context, k string) (i CameraItem, ok bool, err error) {
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

func (c *TiersCache) Set(ctx context.Context, item CameraItem) error {
	return c.innerSet(item)
}

func (c *TiersCache) innerGet(k string) (CameraItem, bool, error) {
	a, ok := c.m.Load(k)
	if ok {
		return a.(CameraItem), ok, nil
	}

	return CameraItem{}, false, nil
}

func (c *TiersCache) innerSet(item CameraItem) error {
	c.m.Store(item.UUID, item)
	return nil
}
