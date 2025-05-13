package holo

import (
	"context"

	"github.com/imroc/req/v3"
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
	isClose bool
	address string
	client  *req.Client

	SN       string
	UUID     string
	DeviceID string
	User     string
	Pwd      string
}

func OpenDevice(addr, username, password string) (*Device, error) {
	c := req.C().
		EnableInsecureSkipVerify().
		SetCommonRetryCount(3).
		DisableCompression().
		SetCommonDigestAuth(username, password)

	return &Device{
		client:  c,
		address: addr,
		User:    username,
		Pwd:     password,
	}, nil
}

func (h *Device) EnableDebug() {
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

func (h *Device) GetMetadataSubscription(ctx context.Context) (data *Subscripions, err error) {
	_, err = h.client.R().
		SetSuccessResult(&data).
		Get(cameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription"))
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
	_, err = h.client.R().
		SetContext(ctx).
		SetSuccessResult(&idList).
		Get(cameraURL(h.address, "/SDCAPI/V1.0/Rest/DeviceID"))
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
	info = new(SysBaseInfo)
	_, err = h.client.R().
		SetSuccessResult(&info).
		Get(cameraURL(h.address, "/SDCAPI/V1.0/MiscIaas/System"))
	return
}
