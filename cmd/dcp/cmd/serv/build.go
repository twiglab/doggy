package serv

import (
	"context"

	"github.com/spf13/viper"
)

type ctxKey struct {
	name string
}

var (
	keyKvHandle = ctxKey{"_kvh_"}
	keyBackend  = ctxKey{"_backend_"}
	keyRootLog  = ctxKey{"_root_log_"}
	keyReg      = ctxKey{"_auto_reg_"}
)

func buildAll(box context.Context, v *viper.Viper) context.Context {
	_, box = buildRootLog(box, v)
	_, box = buildKVHandle(box, v)
	_, box = buildReg(box, v)
	_, box = buildBackend2(box, v)
	return box
}
