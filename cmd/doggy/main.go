package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/doggy"
)

func main() {
	h := &doggy.HoleHandl{}

	mux := chi.NewMux()

	mux.Put("/SDCEntry", doggy.DeviceRegisterUpload(h))
	mux.Post("/SDCEntry", doggy.MetadataEntry(h))
	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("*******")
		fmt.Println(r.Method)
		fmt.Println(r.URL)
		fmt.Println("*******")
	})

	if err := http.ListenAndServe(":10005", mux); err != nil {
		log.Fatal(err)
	}

}
