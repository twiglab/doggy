package pf

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PlatformHandle(h *Handle) http.Handler {
	r := chi.NewRouter()

	r.Put("/SDCEntry", DeviceAutoRegisterUpload(h))
	r.Put("/nat", DeviceAutoRegisterUpload(h))
	r.Put("/1", DeviceAutoRegisterUpload(h))

	r.Post("/SDCEntry", MetadataEntryUpload(h))
	r.Post("/MetadataEntry", MetadataEntryUpload(h))
	r.Post("/upload", MetadataEntryUpload(h))
	r.Post("/1", MetadataEntryUpload(h))

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hahha~~, url = %s, meth = %s, ssl = %t", r.URL.String(), r.Method, r.TLS != nil)
	})
	return r
}
