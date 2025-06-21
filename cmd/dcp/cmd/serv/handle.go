package serv

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/out"
	"github.com/twiglab/doggy/page"
	"github.com/twiglab/doggy/pf"
	"github.com/twiglab/doggy/taosdb"
)

func pageHandle(ctx context.Context, _ AppConf) http.Handler {
	loader := ctx.Value(keyEhc).(page.Loader)
	toucher := ctx.Value(keyToucher).(pf.Toucher)
	p := page.NewPage(loader, toucher)
	return page.AdminPage(p)
}

func outHandle(_ context.Context, conf AppConf) http.Handler {
	switch backendName(conf) {
	case bNameTaos:
		db := MustOpenTaosDB(conf)
		return out.OutHandle(out.NewOutServ(&taosdb.OutS{DB: db}))
	}
	return out.OutHandle(out.NewOutServ(&out.UnimplOut{}))
}

func pfHandle(ctx context.Context, conf AppConf) http.Handler {

	cmdb := ctx.Value(keyCmdb).(pf.DeviceResolver)
	ehc := ctx.Value(keyEhc).(pf.Cache)
	toucher := ctx.Value(keyToucher).(pf.Toucher)

	var backups []holo.SubscriptionReq
	for _, b := range conf.SubsConf.Backups {
		backups = append(backups, MustSubReq(holo.SubReq(b)))
	}

	autoSub := &pf.AutoSub{
		DeviceResolver: cmdb,
		CacheSetter:    ehc,
		MainSub:        MustSubReq(holo.SubReq(conf.SubsConf.Main)),
		Backups:        backups,
		Muti:           conf.SubsConf.Muti,
	}

	var backend pfh
	if backendName(conf) != bNameNone {
		backend = ctx.Value(keyBackend).(*taosdb.Schemaless)
	}

	h := pf.NewHandle(
		pf.WithDeviceRegister(autoSub),
		pf.WithCountHandler(backend),
		pf.WithDensityHandler(backend),
		pf.WithToucher(toucher),
		pf.WithCache(pf.NewTieredCache(ehc)),
	)

	return pf.PlatformHandle(h)
}

func FullHandler(ctx context.Context, conf AppConf) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer, middleware.RequestID)
	mux.Mount("/pf", pfHandle(ctx, conf))
	mux.Mount("/admin", pageHandle(ctx, conf))
	mux.Mount("/debug", middleware.Profiler())
	mux.Mount("/jsonrpc", outHandle(ctx, conf))
	return mux
}
