package serv

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/taosdata/driver-go/v3/ws/schemaless"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/pf"
	"github.com/twiglab/doggy/taosdb"
)

const (
	keyEhc     = "_ehc_"
	keyCmdb    = "_cmdb_"
	keyBackend = "_backend_"
	keyRootLog = "_root_log_"
	keyToucher = "_toucher_"

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

func buildEntCache(ctx context.Context, conf AppConf) (pf.Cache, context.Context) {
	eh := orm.NewEntHandle(MustEntClient(conf.DBConf))
	return eh, context.WithValue(ctx, keyEhc, eh)
}

func buildCmdb(ctx context.Context, conf AppConf) (*pf.CameraDB, context.Context) {
	cmdb := &pf.CameraDB{
		User: conf.CameraDBConf.CsvCameraDB.CameraUser,
		Pwd:  conf.CameraDBConf.CsvCameraDB.CameraPwd,
	}

	return cmdb, context.WithValue(ctx, keyCmdb, cmdb)
}

func buildTaos(ctx context.Context, conf AppConf) (*taosdb.Schemaless, context.Context) {
	url := taosdb.SchemalessURL(
		conf.BackendConf.TaosDBConf.Addr,
		conf.BackendConf.TaosDBConf.Port,
	)
	sc := schemaless.NewConfig(url, 1,
		schemaless.SetDb(conf.BackendConf.TaosDBConf.DBName),
		schemaless.SetAutoReconnect(true),
		schemaless.SetUser(conf.BackendConf.TaosDBConf.Username),
		schemaless.SetPassword(conf.BackendConf.TaosDBConf.Password),
		schemaless.SetReadTimeout(5*time.Second),
		schemaless.SetWriteTimeout(5*time.Second),
	)
	s, err := schemaless.NewSchemaless(sc)
	if err != nil {
		log.Fatal(err)
	}

	schema := taosdb.NewSchLe(s)
	return schema, context.WithValue(ctx, keyBackend, schema)
}

func buildBackend(ctx context.Context, conf AppConf) (pfh, context.Context) {
	switch backendName(conf) {
	case bNameTaos:
		return buildTaos(ctx, conf)
	}
	return nil, ctx
}

func buildToucher(ctx context.Context, _ AppConf) (*pf.InMomoryTouch, context.Context) {
	t := &pf.InMomoryTouch{}
	return t, context.WithValue(ctx, keyToucher, t)
}

func buildAll(box context.Context, conf AppConf) context.Context {
	_, box = buildRootlogger(box, conf)
	_, box = buildEntCache(box, conf)
	_, box = buildCmdb(box, conf)
	_, box = buildBackend(box, conf)
	_, box = buildToucher(box, conf)
	return box
}
