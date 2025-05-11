package serv

import (
	"context"
	"net/http"

	"github.com/twiglab/doggy/out"
)

func outHandle(ctx context.Context, conf AppConf) http.Handler {
	return out.OutHandle(out.NewOutServ())
}
