package holo

import (
	"context"
	"time"

	"resty.dev/v3"
)

func deviceURL(addr, path string) string {
	url := "https://" + addr + path
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
	c := resty.NewWithTransportSettings(
		&resty.TransportSettings{
			DialerTimeout:   10 * time.Second,
			IdleConnTimeout: 20 * time.Second,
		}).SetDigestAuth(username, password)

	return &Device{
		client: c,
		addr:   addr,
	}, nil
}

func CloseDevice(d *Device) error {
	return d.Close()
}

func (h *Device) MetadataSubscription(ctx context.Context, req SubscriptionReq) (*CommonResponseID, error) {
	cr := &CommonResponseID{}

	_, err := h.client.R().
		SetBody(req).
		SetResult(cr).
		SetError(cr).
		Post(deviceURL(h.addr, "/SDCAPI/V2.0/Metadata/Subscription"))

	if err != nil {
		return cr, err
	}

	if cr.IsErr() {
		return cr, cr
	}

	return cr, nil
}

func (h *Device) GetMetadataSubscription(ctx context.Context) (*Subscripions, error) {
	data := &Subscripions{}

	_, err := h.client.R().
		SetResult(data).
		Get(deviceURL(h.addr, "/SDCAPI/V2.0/Metadata/Subscription"))

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
