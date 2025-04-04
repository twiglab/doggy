package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/doggy"
)

func main() {
	h := &doggy.HoleHandl{}

	mux := chi.NewMux()

	mux.Put("/SDCEntry", doggy.DeviceRegisterUpload(h))

	if err := http.ListenAndServe(":10005", mux); err != nil {
		log.Fatal(err)
	}

}
