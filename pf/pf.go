package pf

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/hx"
)

func DeviceAutoRegisterUpload(h *Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data holo.DeviceAutoRegisterData
		if err := hx.BindAndClose(r, &data); err != nil {
			_ = hx.JsonTo(http.StatusInternalServerError,
				holo.CommonResponseFailedError(r.URL.Path, err), w)
			return
		}

		ctx := r.Context()

		if err := h.HandleAutoRegister(ctx, data); err != nil {
			slog.ErrorContext(ctx, "handleAutoReg error",
				slog.String("module", "hapi"),
				slog.Any("error", err))

			_ = hx.JsonTo(http.StatusInternalServerError,
				holo.CommonResponseFailedError(r.URL.Path, err), w)
			return
		}

		_ = hx.JsonTo(http.StatusOK, holo.CommonResponseOK(r.URL.Path), w)
	}
}

func MetadataEntryUpload(h *Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data holo.MetadataObjectUpload
		if err := hx.BindAndClose(r, &data); err != nil {
			slog.ErrorContext(r.Context(), "error-01",
				slog.String("method", "MetadataEntryUpload"),
				slog.Any("error", err))

			http.Error(w, "error-01", http.StatusInternalServerError)
			return
		}

		if err := h.HandleMetadata(r.Context(), data); err != nil {
			slog.ErrorContext(r.Context(), "error-02",
				slog.String("method", "MetadataEntryUpload"),
				slog.Any("error", err))

			http.Error(w, "error-02", http.StatusInternalServerError)
			return
		}

		hx.NoContent(w)
	}
}

func PlatformHandle(h *Handle) http.Handler {
	r := chi.NewRouter()
	r.Put("/nat", DeviceAutoRegisterUpload(h))
	r.Put("1", DeviceAutoRegisterUpload(h))
	r.Post("/upload", MetadataEntryUpload(h))
	r.Post("/2", MetadataEntryUpload(h))
	return r
}
