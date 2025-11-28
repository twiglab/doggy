package serv

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/taosdata/driver-go/v3/ws/schemaless"
	"github.com/twiglab/doggy/backend/taosdb"
	"github.com/twiglab/doggy/kv"
	"github.com/twiglab/doggy/pf"
)

const (
	bNameTaos = "taos"
	bNameNone = "none"
)

type ctxKey struct {
	name string
}

var (
	keyKvHandle = ctxKey{"_ehc_"}
	keyCmdb     = ctxKey{"_cmdb_"}
	keyBackend  = ctxKey{"_backend_"}
	keyRootLog  = ctxKey{"_root_log_"}
)

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

func buildKVHandle(ctx context.Context, conf AppConf) (*kv.Handle, context.Context) {
	h := MustOpenKV(conf)
	return h, context.WithValue(ctx, keyKvHandle, h)
}

func buildCmdb(ctx context.Context, conf AppConf) (*pf.CameraDB, context.Context) {
	cmdb := &pf.CameraDB{
		User:     conf.CameraDBConf.CsvCameraDB.CameraUser,
		Pwd:      conf.CameraDBConf.CsvCameraDB.CameraPwd,
		UseHttps: true,
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

func buildBackend(ctx context.Context, conf AppConf) (pf.DataHandler, context.Context) {
	switch backendName(conf) {
	case bNameTaos:
		return buildTaos(ctx, conf)
	}
	return nil, ctx
}

func buildAll(box context.Context, conf AppConf) context.Context {
	_, box = buildRootlogger(box, conf)
	_, box = buildKVHandle(box, conf)
	_, box = buildCmdb(box, conf)
	_, box = buildBackend(box, conf)
	return box
}
