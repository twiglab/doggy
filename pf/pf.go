package pf

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PlatformHandle(h *HoloHandle) http.Handler {
	mux := chi.NewMux()

	mux.Put("/SDCEntry", DeviceAutoRegisterUpload(h))
	mux.Put("/nat", DeviceAutoRegisterUpload(h))

	mux.Post("/SDCEntry", MetadataEntryUpload(h))
	mux.Post("/MetadataEntry", MetadataEntryUpload(h))

	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hahha~~, url = %s, meth = %s, ssl = %t", r.URL.String(), r.Method, r.TLS != nil)
	})

	return mux
}
