package serv

import (
	"context"
	"log"
	"log/slog"

	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/pf"
	"github.com/twiglab/doggy/taosdb"
)

const (
	key_eh     = "_eh"
	keyCmdb    = "_cmdb_"
	keyBackend = "_backend_"
	keyRootLog = "_root_log_"

	bNameTaos = "taos"
	bNameNone = "none"
)

type pfh interface {
	pf.CountHandler
	pf.DensityHandler
}

func backendName(conf AppConf) string {
	switch conf.BackendConf.Use {
	case "taos", "TAOS":
		return bNameTaos
	}
	return bNameNone
}

func buildRootlogger(ctx context.Context, conf AppConf) (*slog.Logger, context.Context) {
	logger := BuildRootLog(conf)
	return logger, context.WithValue(ctx, keyRootLog, logger)
}

func buildEntHandle(ctx context.Context, conf AppConf) (*orm.EntHandle, context.Context) {
	eh := orm.NewEntHandle(MustEntClient(conf.DBConf))
	return eh, context.WithValue(ctx, key_eh, eh)
}

func buildCmdb(ctx context.Context, conf AppConf) (*pf.CsvCameraDB, context.Context) {
	cmdb := pf.NewCsvCameraDB(
		conf.CameraDBConf.CsvCameraDB.CsvFile,
		conf.CameraDBConf.CsvCameraDB.CameraUser,
		conf.CameraDBConf.CsvCameraDB.CameraPwd,
	)

	if err := cmdb.Load(ctx); err != nil {
		log.Fatal(err)
	}

	return cmdb, context.WithValue(ctx, keyCmdb, cmdb)
}

func buildTaos(ctx context.Context, conf AppConf) (*taosdb.Schemaless, context.Context) {
	schema, err := taosdb.NewSchemaless(taosdb.Config{
		Addr:     conf.BackendConf.TaosDBConf.Addr,
		Port:     conf.BackendConf.TaosDBConf.Port,
		Protocol: conf.BackendConf.TaosDBConf.Protocol,
		Username: conf.BackendConf.TaosDBConf.Username,
		Password: conf.BackendConf.TaosDBConf.Password,
		DBName:   conf.BackendConf.TaosDBConf.DBName,
	})
	if err != nil {
		log.Fatal(err)
	}
	return schema, context.WithValue(ctx, keyBackend, schema)
}

func buildBackend(ctx context.Context, conf AppConf) (pfh, context.Context) {
	switch backendName(conf) {
	case bNameTaos:
		return buildTaos(ctx, conf)
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
