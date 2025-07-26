package kv

import (
	"context"
	"time"

	"github.com/twiglab/doggy/pf"
)

type Touch struct {
	H *Handle

	TTL int64
}

func (t *Touch) Get(ctx context.Context, me string) (time.Time, bool, error) {
	return t.H.TouchLast(ctx, me)
}

func (t *Touch) Set(ctx context.Context, me string, now time.Time) error {
	return t.H.TouchChannel(ctx, me, now, t.TTL)
}

type Upload struct {
	H *Handle
}

func (u *Upload) Upload(ctx context.Context, channels []pf.Channel) error {
	return u.H.SetChannels(ctx, channels)
}

type ChannelCache struct {
	H *Handle
}

func (c *ChannelCache) Get(ctx context.Context, key string) (pf.Channel, bool, error) {
	return c.H.GetChannel(ctx, key)
}

func (c *ChannelCache) Set(ctx context.Context, _ string, ch pf.Channel) error {
	return c.H.SetChannel(ctx, ch)
}
