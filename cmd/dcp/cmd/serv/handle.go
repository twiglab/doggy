package serv

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/page"
	"github.com/twiglab/doggy/pf"
)

func pageHandle(ctx context.Context, _ AppConf) http.Handler {
	loader := ctx.Value(keyEhc).(page.Loader)
	toucher := ctx.Value(keyToucher).(pf.Toucher)
	p := page.NewPage(loader, toucher)
	return page.AdminPage(p)
}

func pfHandle(ctx context.Context, conf AppConf) http.Handler {

	cmdb := ctx.Value(keyCmdb).(pf.DeviceResolver)
	ehc := ctx.Value(keyEhc).(*orm.EntHandle)
	toucher := ctx.Value(keyToucher).(pf.Toucher)

	var backups []holo.SubscriptionReq
	for _, b := range conf.SubsConf.Backups {
		backups = append(backups, MustSubReq(holo.SubReq(b)))
	}

	autoSub := &pf.AutoSub{
		DeviceResolver: cmdb,
		Uploader:       ehc,
		MainSub:        MustSubReq(holo.SubReq(conf.SubsConf.Main)),
		Backups:        backups,
		Muti:           conf.SubsConf.Muti,
	}

	cache := pf.NewTiersCache()
	cache.SetSecond(ehc)

	var backend pfh
	if backendName(conf) != bNameNone {
		backend = ctx.Value(keyBackend).(pfh)
	}

	h := pf.NewHandle(
		pf.WithDeviceRegister(autoSub),
		pf.WithCountHandler(backend),
		pf.WithDensityHandler(backend),
		pf.WithToucher(toucher),
		pf.WithCache(cache),
	)

	return pf.PlatformHandle(h)
}

func FullHandler(ctx context.Context, conf AppConf) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Mount("/pf", pfHandle(ctx, conf))
	mux.Mount("/admin", pageHandle(ctx, conf))
	mux.Mount("/debug", middleware.Profiler())
	return mux
}
