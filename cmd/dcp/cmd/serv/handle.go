package serv

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twiglab/doggy/idb"
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
	case bNameIDB:
		p := ctx.Value(keyBackend).(*idb.IdbPoint)
		acc := idb.NewIdbOut(p)
		return out.OutHandle(out.NewOutServ(acc))
	case bNameTaos:
		db := MustOpenTaosDB(conf)
		return out.OutHandle(out.NewOutServ(&taosdb.OutS{DB: db}))
	}
	return out.OutHandle(out.NewOutServ(&out.UnimplOut{}))
}

func pfTestHandle() http.Handler {
	return pf.PlatformHandle(pf.NewHandle())
}

func pfHandle(ctx context.Context, conf AppConf) http.Handler {
	eh := ctx.Value(key_eh).(pf.UploadHandler)
	fixUser := ctx.Value(keyCmdb).(pf.DeviceResolver)

	autoSub := &pf.AutoSub{
		DeviceResolver: fixUser,
		UploadHandler:  eh,

		MetadataURL: conf.AutoRegConf.MetadataURL,
		Addr:        conf.AutoRegConf.Addr,
		Port:        conf.AutoRegConf.Port,
	}
	h := pf.NewHandle(pf.WithDeviceRegister(autoSub))

	if backendName(conf) != bNameNone {
		backend := ctx.Value(keyBackend).(pfh)
		h.SetCountHandler(backend)
		h.SetDensityHandler(backend)
	}

	return pf.PlatformHandle(h)
}

func MainHandler(ctx context.Context, conf AppConf) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Mount("/", pfTestHandle())
	mux.Mount("/pf", pfHandle(ctx, conf))
	mux.Mount("/admin", pageHandle(ctx, conf))

	mux.Mount("/debug", middleware.Profiler())

	mux.Mount("/jsonrpc", outHandle(ctx, conf))

	return mux
}
