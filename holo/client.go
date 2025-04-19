package holo

import (
	"context"
	"time"

	"resty.dev/v3"
)

func cameraURL(addr, path string) string {
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

func (h *Device) PostMetadataSubscription(ctx context.Context, req SubscriptionReq) (*CommonResponseID, error) {
	cr := &CommonResponseID{}

	_, err := h.client.R().
		SetBody(req).
		SetResult(cr).
		SetError(cr).
		Post(cameraURL(h.addr, "/SDCAPI/V2.0/Metadata/Subscription"))

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
		Get(cameraURL(h.addr, "/SDCAPI/V2.0/Metadata/Subscription"))

	if err != nil {
		return data, err
	}

	return data, nil
}

func (h *Device) Reboot(ctx context.Context) (*RebootResp, error) {
	resp := &RebootResp{}

	_, err := h.client.R().
		SetResult(resp).
		SetError(resp).
		Post(cameraURL(h.addr, "/HSAPI/V1/System/Reboot"))

	if err != nil {
		return resp, err
	}

	if resp.IsErr() {
		return resp, resp
	}
	return resp, nil
}

func (h *Device) Close() error {
	if !h.isClose {
		return h.client.Close()
	}
	return nil
}
