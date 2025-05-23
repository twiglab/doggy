package pf

import (
	"log"
	"net"
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

		// fix SDC ver < 9.0.0
		if data.IpAddr == "" {
			if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				data.IpAddr = ip
			}
		}

		if err := h.HandleAutoRegister(r.Context(), data); err != nil {
			_ = hx.JsonTo(http.StatusInternalServerError,
				holo.CommonResponseFailedError(r.URL.Path, err), w)
			return
		}

		log.Println("register device ok, sn = ", data.SerialNumber)
		_ = hx.JsonTo(http.StatusOK, holo.CommonResponseOK(r.URL.Path), w)
	}
}

func MetadataEntryUpload(h *Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data holo.MetadataObjectUpload
		if err := hx.BindAndClose(r, &data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := h.HandleMetadata(r.Context(), data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
