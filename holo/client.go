package holo

import (
	"resty.dev/v3"
)

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
