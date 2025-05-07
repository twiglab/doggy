package holo

import (
	"context"
	"crypto/tls"
	"net/http"
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
	c := resty.NewWithClient(&http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			IdleConnTimeout:     20 * time.Second,
			DisableKeepAlives:   true,
			MaxIdleConns:        3,
			TLSHandshakeTimeout: 5 * time.Second,
			DisableCompression:  true,
		},
	}).SetDigestAuth(username, password)

	return &Device{
		client: c,
		Addr:   addr,
		User:   username,
		Pwd:    password,
	}, nil
}

func (h *Device) EnableDebug() {
	h.client.SetDebug(true)
}

func (h *Device) Close() error {
	if !h.isClose {
		return h.client.Close()
	}
	return nil
}

func (h *Device) PostMetadataSubscription(ctx context.Context, req SubscriptionReq) (resp *CommonResponse, err error) {
	cr := new(CommonError)

	_, err = h.client.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(cr).
		SetError(cr).
		Post(cameraURL(h.Addr, "/SDCAPI/V2.0/Metadata/Subscription"))

	return
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

func (h *Device) Reboot(ctx context.Context) (resp *CommonResponse, err error) {
	resp = new(CommonResponse)
	_, err = h.client.R().
		SetContext(ctx).
		SetResult(resp).
		SetError(resp).
		Post(cameraURL(h.Addr, "/SDCAPI/V1.0/System/Reboot"))
	return
}

func (h *Device) GetDeviceID(ctx context.Context) (idList *DeviceIDList, err error) {
	idList = new(DeviceIDList)
	_, err = h.client.R().
		SetContext(ctx).
		SetResult(&idList).
		Get(cameraURL(h.Addr, "/SDCAPI/V1.0/Rest/DeviceID"))
	return
}

func (h *Device) PutDeviceID(ctx context.Context, idList DeviceIDList) (resp *CommonResponse, err error) {
	resp = new(CommonResponse)
	_, err = h.client.R().
		SetContext(ctx).
		SetBody(idList).
		SetResult(resp).
		SetError(resp).
		Put(cameraURL(h.Addr, "/SDCAPI/V1.0/Rest/DeviceID"))
	return
}

func (h *Device) GetSysBaseInfo(ctx context.Context) (info *SysBaseInfo, err error) {
	info = new(SysBaseInfo)
	_, err = h.client.R().
		SetResult(info).
		Get(cameraURL(h.Addr, "/SDCAPI/V1.0/MiscIaas/System"))
	return
}
