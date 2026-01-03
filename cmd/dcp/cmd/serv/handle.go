package serv

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twiglab/doggy/kv"
	"github.com/twiglab/doggy/page"
	"github.com/twiglab/doggy/pf"
)

func pageHandle(ctx context.Context) http.Handler {
	kvh := ctx.Value(keyKvHandle).(*kv.Handle)
	toucher := kv.NewTouch(kvh, 90)
	p := page.NewPage(kvh, toucher)
	return page.AdminPage(p)
}

func pfHandle(ctx context.Context, project string) http.Handler {
	autoSub := ctx.Value(keyReg).(pf.DeviceRegister)
	kvh := ctx.Value(keyKvHandle).(*kv.Handle)
	backend := ctx.Value(keyBackend).(pf.DataHandler)

	cache := pf.NewTiersCache[string, pf.ChannelExtra]().WithSecond(&kv.ChannelCache{H: kvh})
	toucher := kv.NewTouch(kvh, 90)

	h := pf.NewMainHandle(
		project,

		pf.WithDeviceRegister(autoSub),
		pf.WithDataHandler(backend),
		pf.WithToucher(toucher),
		pf.WithCache(cache),
	)

	return pf.PlatformHandle(h)
}

func MainHandle(ctx context.Context, project string) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Mount("/pf", pfHandle(ctx, project))
	mux.Mount("/admin", pageHandle(ctx))
	mux.Mount("/debug", middleware.Profiler())
	return mux
}
