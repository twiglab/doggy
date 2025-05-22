package serv

import (
	"context"
	"net/http"

	"github.com/twiglab/doggy/idb"
	"github.com/twiglab/doggy/pf"
)

func pfHandle(ctx context.Context, conf AppConf) http.Handler {
	eh := ctx.Value(key_eh).(pf.UploadHandler)
	fixUser := ctx.Value(key_resolve).(pf.DeviceResolver)

	autoSub := &pf.AutoSub{
		DeviceResolver: fixUser,
		UploadHandler:  eh,

		MetadataURL: conf.AutoRegConf.MetadataURL,
		Addr:        conf.AutoRegConf.Addr,
		Port:        conf.AutoRegConf.Port,
	}
	h := pf.NewHandle(pf.WithDeviceRegister(autoSub))

	if !conf.BackendConf.NoBackend {
		idb3 := ctx.Value(key_idb3).(*idb.IdbPoint)
		h.SetCountHandler(idb3)
		h.SetDensityHandler(idb3)
	}
	return pf.PlatformHandle(h)
}

func pfTestHandle() http.Handler {
	return pf.PlatformHandle(pf.NewHandle())
}
