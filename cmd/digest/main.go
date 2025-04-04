package main

import (
	"log"

	"github.com/icholy/digest"
	"github.com/twiglab/doggy/holo"
	"resty.dev/v3"
)

type X struct {
	Enable  int
	PeerUrl string
}

func main() {
	var cr holo.CommonResponse

	c := resty.New()
	c.SetTransport(&digest.Transport{
		Username: "",
		Password: "",
	})

	_, err := c.R().SetBody(&X{Enable: 1, PeerUrl: ""}).
		SetResult(&cr).
		Post("")

	if err != nil {
		log.Fatal(err)
	}

	log.Println(cr)
}
