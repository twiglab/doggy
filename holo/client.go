package holo

import (
	"context"

	"resty.dev/v3"
)

func url(addr, path string) string {
	return "https://" + addr + path
}

const (
// metadata_subscription = "/SDCAPI/V2.0/Metadata/Subscription"
)

type PlatformConfig struct {
	BaseURL string
	Address string
	Port    int

	DeviceCommonConfig DeviceCommonConfig
}

type DeviceCommonConfig struct {
	Username string
	Password string
}

type Device struct {
	client *resty.Client

	addr string

	config PlatformConfig
}

func NewDevice(addr string, conf PlatformConfig) *Device {
	c := resty.New().SetDigestAuth(conf.DeviceCommonConfig.Username, conf.DeviceCommonConfig.Password)
	return &Device{
		client: c,
		addr:   addr,
		config: conf,
	}
}

func (h *Device) MetadataSubscription(ctx context.Context, req MetadataSubscriptionReq) (*CommonResponseID, error) {
	cr := &CommonResponseID{}

	_, err := h.client.R().
		SetBody(req).
		SetResult(cr).
		SetError(cr).
		Post(url(h.addr, "/SDCAPI/V2.0/Metadata/Subscription"))

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
		Get(url(h.addr, "/SDCAPI/V2.0/Metadata/Subscription"))

	if err != nil {
		return data, err
	}

	return data, nil
}

func (h *Device) Close() error {
	return h.client.Close()
}
