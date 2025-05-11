package serv

import (
	"context"
	"net/http"

	"github.com/twiglab/doggy/page"
	"github.com/twiglab/doggy/pf"
)

func pageHandle(ctx context.Context, conf AppConf) http.Handler {
	loader := ctx.Value(key_eh).(pf.DeviceLoader)
	p := page.NewPage(loader)
	return page.AdminPage(p)
}
