package doggy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/doggy/holo"
)

func JsonTo(code int, resp any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	return enc.Encode(resp)
}

func Bind(r *http.Request, p any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(p); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func HumanCountUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data holo.HumanCountUploadData
		if err := Bind(r, &data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeviceAutoRegisterUpload(h *HoloHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.URL)
		fmt.Println(r.Method)
		fmt.Println(r.TLS != nil)

		var data holo.DeviceAutoRegisterData
		if err := Bind(r, &data); err != nil {
			_ = JsonTo(http.StatusInternalServerError, &holo.CommonResponse{
				RequestUrl:   "/SDCEntry",
				StatusCode:   -1,
				StatusString: err.Error(),
			}, w)
			return
		}

		if err := h.HandleAutoRegister(r.Context(), data); err != nil {
			_ = JsonTo(http.StatusInternalServerError, &holo.CommonResponse{
				RequestUrl:   "/SDCEntry",
				StatusCode:   -1,
				StatusString: err.Error(),
			}, w)
			return
		}

		_ = JsonTo(http.StatusInternalServerError, &holo.CommonResponse{
			RequestUrl:   "/SDCEntry",
			StatusCode:   0,
			StatusString: "OK",
		}, w)
	}
}

func MetadataEntryUpload(h *HoloHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("*******")
		fmt.Println("xxxxxxxx")
		fmt.Println(r.Method)
		fmt.Println(r.URL)
		fmt.Println("*******")

		var data M
		if err := Bind(r, &data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.HandleMeta(r.Context(), data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}

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
		fmt.Println("*******")
		fmt.Println(r.Method)
		fmt.Println(r.URL)
		fmt.Println(r.TLS != nil)
		fmt.Println("*******")
		fmt.Fprintf(w, "url = %s, meth = %s, ssl = %t", r.URL.String(), r.Method, r.TLS != nil)
	})

	return mux
}
