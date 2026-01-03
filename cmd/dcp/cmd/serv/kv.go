package serv

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"github.com/twiglab/doggy/kv"
)

func buildKVHandle(ctx context.Context, v *viper.Viper) (*kv.Handle, context.Context) {
	h := mustOpenKV(v)
	return h, context.WithValue(ctx, keyKvHandle, h)
}

func mustOpenKV(v *viper.Viper) *kv.Handle {
	urls := v.GetStringSlice("etcd.urls")
	h, err := kv.FromURLs(urls)
	if err != nil {
		log.Fatal(err)
	}
	return h
}
