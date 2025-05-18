package holo

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/icholy/digest"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Device struct {
	doer Doer

	address string

	SN       string
	UUID     string
	DeviceID string
	User     string
	Pwd      string
}

func OpenDevice(addr, user, pwd string) (*Device, error) {
	c := &http.Client{
		Transport: &digest.Transport{
			Username: user,
			Password: pwd,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	return &Device{
		doer:    c,
		address: addr,
		User:    user,
		Pwd:     pwd,
	}, nil
}

func (h *Device) PostMetadataSubscription(ctx context.Context, sub SubscriptionReq) (*CommonResponse, error) {
	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	if err := enc.Encode(sub); err != nil {
		return nil, err
	}

	url := cameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}

	resp, err := h.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewApiError(resp.StatusCode, resp.Status)
	}

	comm := &CommonResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(comm); err != nil {
		return nil, err
	}
	return comm, nil
}

func (h *Device) DeleteMetadataSubscription(ctx context.Context) (*CommonResponse, error) {
	url := cameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription")
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewApiError(resp.StatusCode, resp.Status)
	}

	comm := &CommonResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(comm); err != nil {
		return nil, err
	}
	return comm, nil
}

func (h *Device) GetMetadataSubscription(ctx context.Context) (*Subscripions, error) {
	url := cameraURL(h.address, "/SDCAPI/V2.0/Metadata/Subscription")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewApiError(resp.StatusCode, resp.Status)
	}

	sub := &Subscripions{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (h *Device) Reboot(ctx context.Context) (*CommonResponse, error) {
	url := cameraURL(h.address, "/SDCAPI/V1.0/System/Reboot")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewApiError(resp.StatusCode, resp.Status)
	}

	comm := &CommonResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(comm); err != nil {
		return nil, err
	}
	return comm, nil
}

func (h *Device) GetDeviceID(ctx context.Context) (*DeviceIDList, error) {
	url := cameraURL(h.address, "/SDCAPI/V1.0/Rest/DeviceID")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewApiError(resp.StatusCode, resp.Status)
	}

	r := &DeviceIDList{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (h *Device) PutDeviceID(ctx context.Context, idList DeviceIDList) (*CommonResponse, error) {
	url := cameraURL(h.address, "/SDCAPI/V1.0/Rest/DeviceID")

	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	if err := enc.Encode(idList); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, buf)
	if err != nil {
		return nil, err
	}

	resp, err := h.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewApiError(resp.StatusCode, resp.Status)
	}

	comm := &CommonResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(comm); err != nil {
		return nil, err
	}
	return comm, nil
}

func (h *Device) GetSysBaseInfo(ctx context.Context) (info *SysBaseInfo, err error) {
	url := cameraURL(h.address, "/SDCAPI/V1.0/MiscIaas/System")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewApiError(resp.StatusCode, resp.Status)
	}

	r := &SysBaseInfo{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (h *Device) Close() error {
	return nil
}
