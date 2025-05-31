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
	loader := ctx.Value(key_eh).(pf.DeviceLoader)
	p := page.NewPage(loader)
	return page.AdminPage(p)
}

func outHandle(ctx context.Context, conf AppConf) http.Handler {
	switch backendName(conf) {
	case bNameTaos:
		db := MustOpenTaosDB(conf)
		return out.OutHandle(out.NewOutServ(&taosdb.OutS{DB: db}))
	}
	return out.OutHandle(out.NewOutServ(&out.UnimplOut{}))
}

func pfTestHandle() http.Handler {
	return pf.PlatformHandle(pf.NewHandle())
}

func pfHandle2(ctx context.Context, conf AppConf) http.Handler {
	uh := ctx.Value(key_eh).(pf.UploadHandler)
	resolver := ctx.Value(keyCmdb).(pf.DeviceResolver)

	autoReg := &pf.AutoReg{
		DeviceResolver: resolver,
		UploadHandler:  uh,
	}
	h := pf.NewHandle(pf.WithDeviceRegister(autoReg))

	if backendName(conf) != bNameNone {
		backend := ctx.Value(keyBackend).(pfh)
		h.SetCountHandler(backend)
		h.SetDensityHandler(backend)
	}

	return pf.PlatformHandle(h)
}

func pfHandle(ctx context.Context, conf AppConf) http.Handler {
	uh := ctx.Value(key_eh).(pf.UploadHandler)
	resolver := ctx.Value(keyCmdb).(pf.DeviceResolver)

	var backups []holo.SubscriptionReq
	for _, b := range conf.SubsConf.Backups {
		backups = append(backups, holo.SubscriptionReq{
			Address:     b.Addr,
			Port:        b.Port,
			MetadataURL: b.MetadataURL,
			TimeOut:     b.TimeOut,
			HttpsEnable: b.HttpsEnable,
		})
	}

	autoSub := &pf.AutoSub{
		DeviceResolver: resolver,
		UploadHandler:  uh,

		MainSub: holo.SubscriptionReq{
			Address:     conf.SubsConf.Main.Addr,
			Port:        conf.SubsConf.Main.Port,
			MetadataURL: conf.SubsConf.Main.MetadataURL,
			TimeOut:     conf.SubsConf.Main.TimeOut,
			HttpsEnable: conf.SubsConf.Main.HttpsEnable,
		},
		Backups: backups,
		MutiSub: conf.AutoRegConf.MutiSub,
	}
	h := pf.NewHandle(pf.WithDeviceRegister(autoSub))

	if backendName(conf) != bNameNone {
		backend := ctx.Value(keyBackend).(pfh)
		h.SetCountHandler(backend)
		h.SetDensityHandler(backend)
	}

	return pf.PlatformHandle(h)
}

func pfBackendHandle(ctx context.Context, conf AppConf) http.Handler {
	h := pf.NewHandle()

	if backendName(conf) != bNameNone {
		backend := ctx.Value(keyBackend).(pfh)
		h.SetCountHandler(backend)
		h.SetDensityHandler(backend)
	}

	return pf.PlatformHandle(h)
}

func FullHandler(ctx context.Context, conf AppConf) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Mount("/", pfTestHandle())
	mux.Mount("/pf", pfHandle(ctx, conf))
	mux.Mount("/admin", pageHandle(ctx, conf))

	mux.Mount("/debug", middleware.Profiler())

	mux.Mount("/jsonrpc", outHandle(ctx, conf))

	return mux
}

func BackendHandler(ctx context.Context, conf AppConf) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Mount("/pf", pfBackendHandle(ctx, conf))
	mux.Mount("/debug", middleware.Profiler())
	mux.Mount("/jsonrpc", outHandle(ctx, conf))

	return mux
}
