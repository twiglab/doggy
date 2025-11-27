package kv

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/tinylib/msgp/msgp"
	"github.com/twiglab/doggy/pf"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const prefix = "/doggy"
const channelPrefix = prefix + "/channel/item"
const touchPrefix = prefix + "/channel/touch"

func channelKey(uuid string) string {
	return channelPrefixKey() + uuid
}

func channelPrefixKey() string {
	return channelPrefix + "/"
}

func touchKey(uuid string) string {
	return touchPrefix + "/" + uuid
}

type Handle struct {
	client *clientv3.Client
}

func FromURLs(urls []string) (*Handle, error) {
	client, err := clientv3.NewFromURLs(urls)
	if err != nil {
		return nil, err
	}
	return &Handle{client: client}, nil
}

func New(client *clientv3.Client) *Handle {
	return &Handle{client: client}
}

func (h *Handle) GetChannel(ctx context.Context, key string) (pf.Channel, bool, error) {
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

func (h *Handle) SetChannel(ctx context.Context, u pf.Channel) error {
	var sb strings.Builder
	sb.Grow(512)
	if err := msgp.Encode(&sb, &u); err != nil {
		return err
	}
	_, err := h.client.Put(ctx, channelKey(u.UUID), sb.String())
	return err
}

func (h *Handle) SetChannels(ctx context.Context, us []pf.Channel) error {
	_, err := concurrency.NewSTM(h.client, func(stm concurrency.STM) error {
		for _, item := range us {
			var sb strings.Builder
			sb.Grow(512)
			if err := msgp.Encode(&sb, &item); err != nil {
				return err
			}
			stm.Put(channelKey(item.UUID), sb.String())
		}
		return nil
	})
	return err
}

func (h *Handle) AllChannels(ctx context.Context) ([]pf.Channel, error) {
	var items []pf.Channel

	resp, err := h.client.Get(ctx, channelPrefixKey(), clientv3.WithPrefix())
	if err != nil {
		return items, err
	}

	for _, v := range resp.Kvs {
		var item pf.Channel
		if err := msgp.Decode(bytes.NewReader(v.Value), &item); err == nil { // err == nil
			items = append(items, item)
		}
	}

	return items, nil
}

func (h *Handle) TouchChannel(ctx context.Context, uuid string, now time.Time, ttl int64) error {
	lease := clientv3.NewLease(h.client)

	lresp, err := lease.Grant(ctx, ttl)
	if err != nil {
		return nil
	}

	sec := now.Unix()
	timeStr := strconv.FormatInt(sec, 36)

	kv := clientv3.NewKV(h.client)
	_, err = kv.Put(ctx, touchKey(uuid), timeStr, clientv3.WithLease(lresp.ID))
	return err
}

func (h *Handle) TouchLast(ctx context.Context, uuid string) (time.Time, bool, error) {
	resp, err := h.client.Get(ctx, touchKey(uuid))
	if err != nil {
		return time.Unix(0, 0), false, err
	}
	if resp.Count <= 0 {
		return time.Unix(0, 0), false, nil
	}

	i, err := strconv.ParseInt(string(resp.Kvs[0].Value), 36, 64)
	if err != nil {
		return time.Unix(0, 0), false, err
	}
	return time.Unix(i, 0), true, nil
}
