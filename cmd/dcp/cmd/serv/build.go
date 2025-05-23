package serv

import (
	"context"
	"log"
	"log/slog"

	"github.com/twiglab/doggy/idb"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/pf"
	"github.com/twiglab/doggy/taosdb"
)

const (
	key_eh          = "_eh"
	key_resolve     = "_resolve"
	key_idb3        = "_idb3"
	key_root_logger = "_root_logger"
)

type pfh interface {
	pf.CountHandler
	pf.DensityHandler
}

func hasBackend(bName string) bool {
	return bName != "" && bName != "none" && bName != "NONE"
}

func buildRootlogger(ctx context.Context, conf AppConf) (*slog.Logger, context.Context) {
	logger := BuildRootLog(conf)
	return logger, context.WithValue(ctx, key_root_logger, logger)
}

func buildEntHandle(ctx context.Context, conf AppConf) (*orm.EntHandle, context.Context) {
	eh := orm.NewEntHandle(MustEntClient(conf.DBConf))
	return eh, context.WithValue(ctx, key_eh, eh)
}

func buildCmdb(ctx context.Context, conf AppConf) (*pf.CsvCameraDB, context.Context) {
	fixUser := pf.NewCsvCameraDB(
		conf.CameraDBConf.CsvCameraDB.CsvFile,
		conf.CameraDBConf.CsvCameraDB.CameraUser,
		conf.CameraDBConf.CsvCameraDB.CameraPwd,
	)
	if err := fixUser.Load(ctx); err != nil {
		log.Fatal(err)
	}
	return fixUser, context.WithValue(ctx, key_resolve, fixUser)
}

func buildIdb3(ctx context.Context, conf AppConf) (*idb.IdbPoint, context.Context) {
	idb3 := idb.NewIdbPoint(MustIdb(conf.BackendConf.InfluxDBConf))
	return idb3, context.WithValue(ctx, key_idb3, idb3)
}

func buildTaos(ctx context.Context, conf AppConf) (*idb.IdbPoint, context.Context) {
	_, err := taosdb.OpenDB(conf.BackendConf.TaosDBConf.Name, conf.BackendConf.TaosDBConf.DSN)
	if err != nil {
		log.Fatal(err)
	}
	return nil, ctx
}

func buildBackend(ctx context.Context, conf AppConf) (pfh, context.Context) {
	switch conf.BackendConf.Use {
	case "taos", "TAOS":
		return buildTaos(ctx, conf)
	case "idb", "influxdb", "influx", "idb3":
		return buildIdb3(ctx, conf)
	}
	return nil, ctx
}

func build(box context.Context, conf AppConf) context.Context {
	_, box = buildRootlogger(box, conf)
	_, box = buildEntHandle(box, conf)
	_, box = buildCmdb(box, conf)
	_, box = buildBackend(box, conf)

	return box
}
