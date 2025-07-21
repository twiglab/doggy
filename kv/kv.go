package kv

import (
	"bytes"
	"context"
	"strings"


	"github.com/tinylib/msgp/msgp"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const prefix = "/doggy"
const channelPrefix = "/doggy/channel/item"
const touchPrefix = "/doggy/channel/touch"

func channelKey(uuid string) string {
	return channelPrefix + "/" + uuid
}

func channelPrefixKey() string {
	return channelPrefix + "/"
}

func touchKey(uuid string) string {
	return touchPrefix + "/" + uuid
}

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

func (h *KeyValHandle) GetChannel(ctx context.Context, key string) (pf.Channel, bool, error) {
	var item pf.Channel

	resp, err := h.client.Get(ctx, channelKey(key))
	if err != nil {
		return item, false, err
	}

	if resp.Count <= 0 {
		return item, false, nil
	}

	if err := msgp.Decode(bytes.NewReader(resp.Kvs[0].Value), &item); err != nil {
		return item, false, err
	}

	return item, true, nil
}

func (h *KeyValHandle) SetChannel(ctx context.Context, u pf.Channel) error {
	var sb strings.Builder
	if err := msgp.Encode(&sb, &u); err != nil {
		return err
	}
	_, err := h.client.Put(ctx, channelKey(u.UUID), sb.String())
	return err
}

func (h *KeyValHandle) SetChannels(ctx context.Context, us []pf.Channel) error {
	_, err := concurrency.NewSTM(h.client, func(stm concurrency.STM) error {
		for _, item := range us {
			var sb strings.Builder
			if err := msgp.Encode(&sb, &item); err != nil {
				return err
			}
			stm.Put(channelKey(item.UUID), sb.String())
		}
		return nil
	})
	return err
}

func (h *KeyValHandle) AllChannels(ctx context.Context) ([]pf.Channel, error) {
	var items []pf.Channel

	resp, err := h.client.Get(ctx, channelPrefixKey(), clientv3.WithFromKey())
	if err != nil {
		return items, err
	}

	for _, v := range resp.Kvs {
		var item pf.Channel
		if err := msgp.Decode(bytes.NewReader(v.Value), &item); err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (h *KeyValHandle) TouchChannel(ctx context.Context, uuid string, ttl int64) error {
	lease := clientv3.NewLease(h.client)

	lresp, err := lease.Grant(ctx, ttl)
	if err != nil {
		return nil
	}

	kv := clientv3.NewKV(h.client)
	_, err = kv.Put(ctx, touchKey(uuid), "1", clientv3.WithLease(lresp.ID))
	return err
}
