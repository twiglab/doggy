package serv

import (
	"context"

	"github.com/twiglab/doggy/pf"
)

type ctxKey struct {
	name string
}

var (
	keyKvHandle = ctxKey{"_kvh_"}
	keyCmdb     = ctxKey{"_cmdb_"}
	keyBackend  = ctxKey{"_backend_"}
	keyRootLog  = ctxKey{"_root_log_"}
)

func buildCmdb(ctx context.Context, conf AppConf) (*pf.CameraDB, context.Context) {
	cmdb := &pf.CameraDB{
		User:     conf.CameraDBConf.CsvCameraDB.CameraUser,
		Pwd:      conf.CameraDBConf.CsvCameraDB.CameraPwd,
		UseHttps: true,
	}

	return cmdb, context.WithValue(ctx, keyCmdb, cmdb)
}

func buildAll(box context.Context, conf AppConf) context.Context {
	_, box = buildRootlogger(box, conf)
	_, box = buildKVHandle(box, conf)
	_, box = buildCmdb(box, conf)
	_, box = buildBackend2(box, conf)
	return box
}
