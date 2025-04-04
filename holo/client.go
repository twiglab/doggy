package holo

import (
	"fmt"

	"resty.dev/v3"
)

type CommonResponse struct {
	RequestUrl   string
	StatusCode   int
	StatusString string
}

func (r *CommonResponse) Error() string {
	return fmt.Sprintf("url = %s, code = %d, str = %s", r.RequestUrl, r.StatusCode, r.StatusString)
}

func (r *CommonResponse) IsErr() bool {
	// 5.2 响应码
	return r.StatusCode != 0
}

type HoloSens struct {
	client *resty.Client
}

func (h *HoloSens) MetadataSubscripton(req MetaSubReq) error {
	cr := &CommonResponse{}

	_, err := h.client.R().
		SetBody(req).
		SetResult(cr).
		SetError(cr).
		Post("")

	if err != nil {
		return err
	}

	if cr.IsErr() {
		return cr
	}

	return nil
}

func (h *HoloSens) Close() error {
	return h.client.Close()
}
