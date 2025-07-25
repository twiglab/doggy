package kv

import (
	"context"
	"time"

	"github.com/twiglab/doggy/pf"
)

type Touch struct {
	h *KeyValHandle

	ttl int64
}

func (t *Touch) Get(ctx context.Context, me string) (time.Time, bool, error) {
	return t.h.TouchLast(ctx, me)
}

func (t *Touch) Set(ctx context.Context, me string, now time.Time) error {
	return t.h.TouchChannel(ctx, me, now, t.ttl)
}

type Upload struct {
	h *KeyValHandle
}

func (u *Upload) Upload(ctx context.Context, channels []pf.Channel) error {
	return u.h.SetChannels(ctx, channels)
}

type ChannelCache struct {
	h *KeyValHandle
}

func (c *ChannelCache) Get(ctx context.Context, key string) (pf.Channel, bool, error) {
	return c.h.GetChannel(ctx, key)
}

func (c *ChannelCache) Set(ctx context.Context, _ string, ch pf.Channel) error {
	return c.h.SetChannel(ctx, ch)
}
