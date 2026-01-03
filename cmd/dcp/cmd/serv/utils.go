package serv

import (
	"log"

	"github.com/twiglab/doggy/holo"
)

func MustSubReq(req holo.SubscriptionReq, err error) holo.SubscriptionReq {
	if err != nil {
		log.Fatal(err)
	}
	return req
}
