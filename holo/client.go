package holo

import (
	"context"

	"github.com/imroc/req/v3"
)

type Device struct {
	client *req.Client

	address string

	User string
	Pwd  string

	useHttps bool
}

func OpenDevice(addr, username, password string, useHttps bool) (*Device, error) {
	c := req.C().
		SetUserAgent("doggy client").
		EnableInsecureSkipVerify().
		DisableCompression().
		SetCommonDigestAuth(username, password)

	return &Device{
		client:  c,
		address: addr,

		User: username,
		Pwd:  password,

		useHttps: useHttps,
	}, nil
}

func ConnectDevice(addr, username, password string) (*Device, error) {
	return OpenDevice(addr, username, password, true)
}

func (h *Device) EnableDebug() {
	h.client.EnableDumpAll()
}

func (h *Device) Close() error {
	return nil
}

func (h *Device) PostMetadataSubscription(ctx context.Context, req SubscriptionReq) (resp *CommonResponse, err error) {
	_, err = h.client.R().
		SetContext(ctx).
		SetBody(req).
		SetSuccessResult(&resp).
		SetErrorResult(&resp).
		Post(CameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription", h.useHttps))
	return
}

func (h *Device) DeleteMetadataSubscription(ctx context.Context) (resp *CommonResponse, err error) {
	_, err = h.client.R().
		SetContext(ctx).
		SetSuccessResult(&resp).
		SetErrorResult(&resp).
		Delete(CameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription", h.useHttps))
	return
}

func (h *Device) GetMetadataSubscription(ctx context.Context) (data *Subscripions, err error) {
	var resp *req.Response
	resp, err = h.client.R().
		SetSuccessResult(&data).
		Get(CameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription", h.useHttps))

	if resp.IsErrorState() {
		err = NewApiError(resp.StatusCode, resp.Status)
	}
	return
}

func (h *Device) Reboot(ctx context.Context) (resp *CommonResponse, err error) {
	_, err = h.client.R().
		SetContext(ctx).
		SetSuccessResult(&resp).
		SetErrorResult(&resp).
		Post(CameraURL(h.address, "/SDCAPI/V1.0/System/Reboot", h.useHttps))
	return
}

func (h *Device) GetDeviceID(ctx context.Context) (idList *DeviceIDList, err error) {
	var resp *req.Response
	resp, err = h.client.R().
		SetContext(ctx).
		SetSuccessResult(&idList).
		Get(CameraURL(h.address, "/SDCAPI/V1.0/Rest/DeviceID", h.useHttps))

	if resp.IsErrorState() {
		err = NewApiError(resp.StatusCode, resp.Status)
	}
	return
}

func (h *Device) PutDeviceID(ctx context.Context, idList DeviceIDList) (resp *CommonResponse, err error) {
	_, err = h.client.R().
		SetContext(ctx).
		SetBody(idList).
		SetSuccessResult(&resp).
		SetErrorResult(&resp).
		Put(CameraURL(h.address, "/SDCAPI/V1.0/Rest/DeviceID", h.useHttps))
	return
}

func (h *Device) GetSysBaseInfo(ctx context.Context) (info *SysBaseInfo, err error) {
	_, err = h.client.R().
		SetSuccessResult(&info).
		Get(CameraURL(h.address, "/SDCAPI/V1.0/MiscIaas/System", h.useHttps))
	return
}
