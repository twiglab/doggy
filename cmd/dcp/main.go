package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twiglab/doggy"
)

func main() {
	h := &doggy.HoloHandle{}

	mux := chi.NewMux()
	mux.Use(middleware.Logger, middleware.Recoverer)

	mux.Put("/SDCEntry", doggy.DeviceAutoRegisterUpload(h))

	mux.Post("/SDCEntry", doggy.MetadataEntryUpload(h))
	mux.Post("/MetadataEntry", doggy.MetadataEntryUpload(h))

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("*******")
		fmt.Println(r.Method)
		fmt.Println(r.URL)
		fmt.Println(r.TLS != nil)
		fmt.Println("*******")
		fmt.Fprintf(w, "url = %s, meth = %s, ssl = %t", r.URL.String(), r.Method, r.TLS != nil)
	})

	/*
		if err := http.ListenAndServeTLS(":10005", "./cert/server.crt", "./cert/server.key", mux); err != nil {
			log.Fatal(err)
		}
	*/

	if err := http.ListenAndServe(":10005", mux); err != nil {
		log.Fatal(err)
	}

}
