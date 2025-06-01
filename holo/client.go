package holo

import (
	"context"

	"github.com/imroc/req/v3"
)

type Device struct {
	address string
	client  *req.Client

	SN   string
	User string
	Pwd  string
}

func OpenDevice(sn, addr, username, password string) (*Device, error) {
	c := req.C().
		SetUserAgent("doggy client").
		EnableInsecureSkipVerify().
		DisableCompression().
		SetCommonDigestAuth(username, password)

	return &Device{
		client:  c,
		address: addr,

		SN:   sn,
		User: username,
		Pwd:  password,
	}, nil
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
		Post(cameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription"))
	return
}

func (h *Device) DeleteMetadataSubscription(ctx context.Context) (resp *CommonResponse, err error) {
	_, err = h.client.R().
		SetContext(ctx).
		SetSuccessResult(&resp).
		SetErrorResult(&resp).
		Delete(cameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription"))
	return
}

func (h *Device) GetMetadataSubscription(ctx context.Context) (data *Subscripions, err error) {
	var resp *req.Response
	resp, err = h.client.R().
		SetSuccessResult(&data).
		Get(cameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription"))

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
		Post(cameraURL(h.address, "/SDCAPI/V1.0/System/Reboot"))
	return
}

func (h *Device) GetDeviceID(ctx context.Context) (idList *DeviceIDList, err error) {
	var resp *req.Response
	resp, err = h.client.R().
		SetContext(ctx).
		SetSuccessResult(&idList).
		Get(cameraURL(h.address, "/SDCAPI/V1.0/Rest/DeviceID"))

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
		Put(cameraURL(h.address, "/SDCAPI/V1.0/Rest/DeviceID"))
	return
}

func (h *Device) GetSysBaseInfo(ctx context.Context) (info *SysBaseInfo, err error) {
	_, err = h.client.R().
		SetSuccessResult(&info).
		Get(cameraURL(h.address, "/SDCAPI/V1.0/MiscIaas/System"))
	return
}
