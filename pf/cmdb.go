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

type innerMemoryCache struct {
	m sync.Map
}

func (i *innerMemoryCache) Get(ctx context.Context, k string) (CameraItem, bool, error) {
	a, ok := i.m.Load(k)
	if ok {
		return a.(CameraItem), ok, nil
	}

	return CameraItem{}, false, nil
}

func (i *innerMemoryCache) Set(ctx context.Context, item CameraItem) error {
	i.m.Store(item.UUID, item)
	return nil
}

type emptyCache string

func (i emptyCache) Get(_ context.Context, _ string) (c CameraItem, ok bool, err error) {
	return
}

func (i emptyCache) Set(_ context.Context, _ CameraItem) (err error) {
	return
}

type TiersCache struct {
	inner  *innerMemoryCache
	second Cache
}

func NewTiersCache() *TiersCache {
	return &TiersCache{
		inner:  &innerMemoryCache{},
		second: emptyCache("x"),
	}
}

func (c *TiersCache) SetSecond(second Cache) {
	c.second = second
}

func (c *TiersCache) Get(ctx context.Context, k string) (i CameraItem, ok bool, err error) {
	if i, ok, err = c.inner.Get(ctx, k); ok {
		return
	}

	if i, ok, err = c.second.Get(ctx, k); ok {
		err = c.inner.Set(ctx, i)
		slog.DebugContext(ctx, "second cache",
			slog.Any("camera", i),
			slog.Any("setCacheErr", err),
			slog.Bool("ok", ok))
	}

	return
}

func (c *TiersCache) Set(ctx context.Context, item CameraItem) error {
	return c.inner.Set(ctx, item)
}
