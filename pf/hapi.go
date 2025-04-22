package pf

import (
	"net"
	"net/http"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/hx"
)

func DeviceAutoRegisterUpload(h *Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data holo.DeviceAutoRegisterData
		if err := hx.Bind(r, &data); err != nil {
			_ = hx.JsonTo(http.StatusInternalServerError, &holo.CommonResponse{
				RequestUrl:   r.URL.Path,
				StatusCode:   -1,
				StatusString: err.Error(),
			}, w)
			return
		}

		// fix SDC ver < 9.0.0
		if data.IpAddr == "" {
			if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				data.IpAddr = ip
			}
		}

		if err := h.HandleAutoRegister(r.Context(), data); err != nil {
			_ = hx.JsonTo(http.StatusInternalServerError, &holo.CommonResponse{
				RequestUrl:   r.URL.Path,
				StatusCode:   -1,
				StatusString: err.Error(),
			}, w)
			return
		}

		_ = hx.JsonTo(http.StatusOK, &holo.CommonResponse{
			RequestUrl:   r.URL.Path,
			StatusCode:   0,
			StatusString: "OK",
		}, w)
	}
}

func MetadataEntryUpload(h *Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data holo.MetadataObjectUpload
		if err := hx.Bind(r, &data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.HandleMetadata(r.Context(), data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
