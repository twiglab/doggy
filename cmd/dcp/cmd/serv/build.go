package serv

import (
	"context"
	"log"
	"log/slog"

	"github.com/twiglab/doggy/idb"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/pf"
)

const (
	key_eh          = "_eh"
	key_resolve     = "_resolve"
	key_idb3        = "_idb3"
	key_root_logger = "_root_logger"
)

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

func buildIDB(ctx context.Context, conf AppConf) (*idb.IdbPoint, context.Context) {
	var idb3 *idb.IdbPoint

	if !conf.BackendConf.NoBackend {
		idb3 = idb.NewIdbPoint(MustIdb(conf.BackendConf.InfluxDBConf))
		return idb3, context.WithValue(ctx, key_idb3, idb3)
	}

	return nil, ctx
}

func build(box context.Context, conf AppConf) context.Context {
	_, box = buildRootlogger(box, conf)
	_, box = buildEntHandle(box, conf)
	_, box = buildCmdb(box, conf)
	_, box = buildIDB(box, conf)

	return box
}
