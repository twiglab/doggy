package serv

import (
	"log"

	"github.com/spf13/viper"
	"github.com/twiglab/doggy/holo"
)

func MustSubReq(req holo.SubscriptionReq, err error) holo.SubscriptionReq {
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func pid(v *viper.Viper) (id string) {
	if id = v.GetString("project.id"); id == "" {
		id = "0000"
	}
	return
}
