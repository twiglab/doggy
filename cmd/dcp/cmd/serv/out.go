package serv

import (
	"context"
	"net/http"

	"github.com/twiglab/doggy/idb"
	"github.com/twiglab/doggy/out"
)

func outHandle(ctx context.Context, conf AppConf) http.Handler {
	p := ctx.Value(key_idb3).(*idb.IdbPoint)
	acc := idb.NewIdbOut(p)
	return out.OutHandle(out.NewOutServ(acc))
}
