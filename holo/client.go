package holo

import (
	"context"
	"fmt"

	"resty.dev/v3"
)

func clientUrl(addr, path string) string {
	url := "https://" + addr + path
	fmt.Println("client url = " + url)
	return url
}

const (
// metadata_subscription = "/SDCAPI/V2.0/Metadata/Subscription"
)

type PlatformInfo struct {
}

type Device struct {
	client  *resty.Client
	isClose bool

	addr string
}

func OpenDevice(addr, username, password string) (*Device, error) {
	c := resty.New().SetDigestAuth(username, password)
	return &Device{
		client: c,
		addr:   addr,
	}, nil
}

func CloseDevice(d *Device) error {
	return d.Close()
}

func (h *Device) MetadataSubscription(ctx context.Context, req MetadataSubscriptionReq) (*CommonResponseID, error) {
	cr := &CommonResponseID{}

	_, err := h.client.R().
		SetBody(req).
		SetResult(cr).
		SetError(cr).
		Post(clientUrl(h.addr, "/SDCAPI/V2.0/Metadata/Subscription"))

	if err != nil {
		return cr, err
	}

	if cr.IsErr() {
		return cr, cr
	}

	return cr, nil
}

func (h *Device) GetMetadataSubscription(ctx context.Context) (*SubscripionsData, error) {
	data := &SubscripionsData{}

	_, err := h.client.R().
		SetResult(data).
		Get(clientUrl(h.addr, "/SDCAPI/V2.0/Metadata/Subscription"))

	if err != nil {
		return data, err
	}

	return data, nil
}

func (h *Device) Close() error {
	if !h.isClose {
		return h.client.Close()
	}
	return nil
}
