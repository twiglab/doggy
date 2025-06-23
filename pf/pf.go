package pf

import (
	"fmt"
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

		slog.InfoContext(ctx, "receive reg data",
			slog.String("module", "hapi"),
			slog.Any("data", data))

		if err := h.HandleAutoRegister(ctx, data); err != nil {

			slog.ErrorContext(ctx, "handleAutoReg error",
				slog.String("module", "hapi"),
				slog.Any("error", err))

			_ = hx.JsonTo(http.StatusInternalServerError,
				holo.CommonResponseFailedError(r.URL.Path, err), w)
			return
		}

		slog.InfoContext(r.Context(), "register device ok",
			slog.String("module", "hapi"),
			slog.String("sn", data.SerialNumber),
			slog.String("ip", data.IpAddr),
		)

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
	r.Put("/1", DeviceAutoRegisterUpload(h))

	r.Post("/upload", MetadataEntryUpload(h))
	r.Post("/1", MetadataEntryUpload(h))
	r.Post("/2", MetadataEntryUpload(h))

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hahha~~, url = %s, meth = %s, ssl = %t", r.URL.String(), r.Method, r.TLS != nil)
	})
	return r
}
