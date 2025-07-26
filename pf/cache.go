package pf

import (
	"context"
	"sync"
)

type Cache[K comparable, V any] interface {
	Get(context.Context, K) (V, bool, error)
	Set(context.Context, K, V) error
}

type localCache[K comparable, V any] struct {
	local sync.Map
}

func (l *localCache[K, V]) Get(_ context.Context, k K) (v V, ok bool, err error) {
	var a any
	if a, ok = l.local.Load(k); ok {
		v = a.(V)
	}
	return
}

func (l *localCache[K, V]) Set(_ context.Context, k K, v V) error {
	l.local.Store(k, v)
	return nil
}

type emptyCache[K comparable, V any] struct{}

func (e emptyCache[K, V]) Get(_ context.Context, _ K) (val V, ok bool, err error) { return }
func (e emptyCache[K, V]) Set(_ context.Context, _ K, _ V) (err error)            { return }

type TiersCache[K comparable, V any] struct {
	local  Cache[K, V]
	second Cache[K, V]
}

func NewTiersCache[K comparable, V any]() *TiersCache[K, V] {
	return &TiersCache[K, V]{
		local:  &localCache[K, V]{},
		second: emptyCache[K, V]{},
	}
}

func (t *TiersCache[K, V]) WithLocal(local Cache[K, V]) *TiersCache[K, V] {
	t.local = local
	return t
}

func (t *TiersCache[K, V]) WithSecond(second Cache[K, V]) *TiersCache[K, V] {
	t.second = second
	return t
}

func (t *TiersCache[K, V]) Get(ctx context.Context, key K) (v V, ok bool, err error) {
	if v, ok, err = t.local.Get(ctx, key); ok {
		return
	}
	if v, ok, err = t.second.Get(ctx, key); ok {
		err = t.local.Set(ctx, key, v)
	}
	return
}

func (t *TiersCache[K, V]) Set(ctx context.Context, k K, v V) error {
	return t.local.Set(ctx, k, v)
}
