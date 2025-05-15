package serv

import (
	"context"
	"log"
	"net/http"

	"github.com/twiglab/doggy/idb"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/pf"
)

const (
	key_eh      = "_eh"
	key_resolve = "_resolve"
	key_idb3    = "_idb3"
)

func buildCtx(conf AppConf) context.Context {

	box := context.Background()

	eh := orm.NewEntHandle(MustEntClient(conf.DBConf))
	box = context.WithValue(box, key_eh, eh)

	fixUser := pf.NewCsvCameraDB(
		conf.CameraDBConf.CsvCameraDB.CsvFile,
		conf.CameraDBConf.CsvCameraDB.CameraUser,
		conf.CameraDBConf.CsvCameraDB.CameraPwd,
	)
	if err := fixUser.Load(); err != nil {
		log.Fatal(err)
	}
	box = context.WithValue(box, key_resolve, fixUser)

	var idb3 *idb.IdbPoint

	if !conf.BackendConf.NoBackend {
		idb3 = idb.NewIdbPoint(MustIdb(conf.BackendConf.InfluxDBConf))
		box = context.WithValue(box, key_idb3, idb3)
	}

	return box
}

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
