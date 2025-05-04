package holo

import (
	"context"
	"time"

	"resty.dev/v3"
)

const (
	HUMMAN_DENSITY = 12
	HUMMAN_COUNT   = 15
)

func cameraURL(addr, path string) string {
	url := "https://" + addr + path
	return url
}

type Device struct {
	client  *resty.Client
	isClose bool

	Addr string
	User string
	Pwd  string
}

func OpenDevice(addr, username, password string) (*Device, error) {
	c := resty.NewWithTransportSettings(
		&resty.TransportSettings{
			DialerTimeout:   10 * time.Second,
			IdleConnTimeout: 20 * time.Second,
		}).SetDigestAuth(username, password)

	return &Device{
		client: c,
		Addr:   addr,
		User:   username,
		Pwd:    password,
	}, nil
}

func (h *Device) Close() error {
	if !h.isClose {
		return h.client.Close()
	}
	return nil
}

func (h *Device) PostMetadataSubscription(ctx context.Context, req SubscriptionReq) (*CommonResponseID, error) {
	cr := &CommonResponseID{}

	_, err := h.client.R().
		SetBody(req).
		SetResult(cr).
		SetError(cr).
		Post(cameraURL(h.Addr, "/SDCAPI/V2.0/Metadata/Subscription"))

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
		Get(cameraURL(h.Addr, "/SDCAPI/V2.0/Metadata/Subscription"))

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
		Post(cameraURL(h.Addr, "/HSAPI/V1/System/Reboot"))

	if err != nil {
		return resp, err
	}

	if resp.IsErr() {
		return resp, resp
	}
	return resp, nil
}

func (h *Device) GetDeviceID(ctx context.Context) (*DeviceIDList, error) {
	ids := &DeviceIDList{}

	_, err := h.client.R().
		SetResult(ids).
		Get(cameraURL(h.Addr, "/SDCAPI/V1.0/Rest/DeviceID"))

	if err != nil {
		return ids, err
	}
	return ids, nil
}

func (h *Device) PutDeviceID(ctx context.Context, idList DeviceIDList) (*CommonResponse, error) {
	resp := &CommonResponse{}

	_, err := h.client.R().
		SetBody(idList).
		SetResult(resp).
		Put(cameraURL(h.Addr, "/SDCAPI/V1.0/Rest/DeviceID"))

	if err != nil {
		return resp, err
	}

	if resp.IsErr() {
		return resp, resp
	}
	return resp, nil
}
