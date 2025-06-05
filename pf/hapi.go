package pf

import (
	"log/slog"
	"net/http"

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
		// ---------------------------
		if data.DeviceVersion.Software == holo.SDC_11_0_0_SPC300 {
			// D3252 SDC 11.0.0.SPC300 重复注册
			// 硬件问题，软件修复
			// 注册成功后，记录自动注册成功的SN号，如果重复注册，直接返回成功
			if h.isSnOk(data.SerialNumber) {
				_ = hx.JsonTo(http.StatusOK, holo.CommonResponseOK(r.URL.Path), w)
				return
			}
		}
		// ---------------------------

		if err := h.HandleAutoRegister(r.Context(), data); err != nil {
			_ = hx.JsonTo(http.StatusInternalServerError,
				holo.CommonResponseFailedError(r.URL.Path, err), w)
			return
		}

		slog.InfoContext(r.Context(), "register device ok",
			slog.String("module", "hapi"),
			slog.String("sn", data.SerialNumber),
			slog.String("ip", data.IpAddr),
		)

		// ---------------------------
		if data.DeviceVersion.Software == holo.SDC_11_0_0_SPC300 {
			// 注册成功后，记录自动注册成功的SN号，如果重复注册，直接返回成功
			h.setSnOk(data.SerialNumber)
		}
		// ---------------------------
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
