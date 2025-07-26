package pf

import (
	"context"
	"sync"
)

type Cache[K, V any] interface {
	Get(context.Context, K) (V, bool, error)
	Set(context.Context, K, V) error
}

type emptyCache[K, V any] struct{}

func (e emptyCache[K, T]) Get(_ context.Context, _ K) (val T, ok bool, err error) { return }
func (e emptyCache[K, T]) Set(_ context.Context, _ K, _ T) (err error)            { return }

type TiersCache[K, V any] struct {
	m      sync.Map
	second Cache[K, V]
}

func NewTiersCache[K any, V any]() *TiersCache[K, V] {
	return NewWithSecond(emptyCache[K, V]{})
}

func NewWithSecond[K any, V any](second Cache[K, V]) *TiersCache[K, V] {
	return &TiersCache[K, V]{
		second: second,
	}
}

func (t *TiersCache[K, V]) Get(ctx context.Context, key K) (v V, ok bool, err error) {
	if v, ok, err = t.innerG(key); ok {
		return
	}
	if v, ok, err = t.second.Get(ctx, key); ok {
		t.m.Store(key, v)
	}
	return
}

func (t *TiersCache[K, V]) Set(ctx context.Context, k K, v V) error {
	t.m.Store(k, v)
	return nil
}

func (t *TiersCache[K, V]) innerG(key K) (V, bool, error) {
	a, ok := t.m.Load(key)
	if ok {
		return a.(V), true, nil
	}

	var v V
	return v, false, nil
}
