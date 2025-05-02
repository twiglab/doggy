package pf

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func PlatformHandle(h *Handle) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Put("/SDCEntry", DeviceAutoRegisterUpload(h))
	mux.Put("/nat", DeviceAutoRegisterUpload(h))

	mux.Post("/SDCEntry", MetadataEntryUpload(h))
	mux.Post("/MetadataEntry", MetadataEntryUpload(h))

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hahha~~, url = %s, meth = %s, ssl = %t", r.URL.String(), r.Method, r.TLS != nil)
	})

	return mux
}
