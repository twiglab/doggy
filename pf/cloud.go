package pf

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/hx"
	"resty.dev/v3"
)

type DataHandler interface {
	HandleData(ctx context.Context, data UploadeData) error
}

type UploadeData struct {
	Common holo.Common
	Target holo.HumanMix
	Tenant Tenant
}

type CloudAction struct {
	client  *resty.Client
	baseURL string
	key     string
}

func NewCloudAction(key, baseURL string) *CloudAction {
	clinet := resty.New().
		SetBaseURL(baseURL).
		SetAuthToken(key)

	return &CloudAction{
		client:  clinet,
		baseURL: baseURL,
		key:     key,
	}
}

func (c *CloudAction) HandleData(ctx context.Context, data UploadeData) error {
	upload := holo.MetadataObjectUpload{
		MetadataObject: holo.MetadataObject{
			Common:     data.Common,
			TargetList: []holo.HumanMix{data.Target},
		},
	}
	r := c.client.R().
		SetContext(ctx).
		SetBody(upload)
	_, err := r.Post("/cloud")
	return err
}

type CloudHandle struct {
	dataHandler DataHandler
}

func CloudUpload(h *CloudHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data UploadeData
		if err := hx.BindAndClose(r, &data); err != nil {
			slog.ErrorContext(r.Context(), "error-01",
				slog.String("method", "MetadataEntryUpload"),
				slog.Any("error", err))

			http.Error(w, "error-01", http.StatusInternalServerError)
			return
		}

		if err := h.dataHandler.HandleData(r.Context(), data); err != nil {
			slog.ErrorContext(r.Context(), "error-02",
				slog.String("method", "MetadataEntryUpload"),
				slog.Any("error", err))

			http.Error(w, "error-02", http.StatusInternalServerError)
			return
		}

		hx.NoContent(w)
	}
}
