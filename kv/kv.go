package main

import (
	"bytes"
	"context"
	"strings"

	"github.com/tinylib/msgp/msgp"
	"github.com/twiglab/doggy/pf"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type KeyValHandle struct {
	client *clientv3.Client
}

func New(urls []string) (*KeyValHandle, error) {
	client, err := clientv3.NewFromURLs(urls)
	if err != nil {
		return nil, err
	}
	return &KeyValHandle{client: client}, nil
}

func (h *KeyValHandle) Get(ctx context.Context, key string) (pf.CameraItem, bool, error) {
	var item pf.CameraItem

	kv := clientv3.NewKV(h.client)

	resp, err := kv.Get(ctx, key)
	if err != nil {
		return item, false, err
	}

	if err := msgp.Decode(bytes.NewReader(resp.Kvs[0].Value), &item); err != nil {
		return item, false, err
	}

	return item, true, nil
}

func (h *KeyValHandle) Set(ctx context.Context, u pf.CameraItem) error {
	var sb strings.Builder
	sb.Grow(512)
	if err := msgp.Encode(&sb, &u); err != nil {
		return err
	}

	h.client.Put(ctx, u.SN, sb.String())
	return nil
}

func (h *KeyValHandle) GetAll(ctx context.Context) ([]pf.CameraItem, error) {
	return nil, nil
}
