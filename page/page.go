package page

import (
	"context"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/doggy/pf"
)

type Loader interface {
	All(context.Context) ([]pf.CameraUpload, error)
}

type Item struct {
	Upload pf.CameraUpload
	Data   pf.CameraData
	TTL    int64
}

type Page struct {
	tpl    *template.Template
	Loader Loader
	Cmdb   *pf.CsvCameraDB
}

func (v *Page) All(ctx context.Context) ([]Item, error) {
	uploads, err := v.Loader.All(ctx)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, u := range uploads {
		data := v.Cmdb.GetBySn(u.SN)
		ttl := v.Cmdb.GetTTL(data.UUID)
		items = append(items, Item{TTL: ttl, Upload: u, Data: data})
	}

	return items, nil
}

func NewPage(loader Loader, cmdb *pf.CsvCameraDB) *Page {
	return &Page{
		tpl:    template.Must(template.ParseFS(tplFS, "tpl/*.tpl")),
		Loader: loader,
		Cmdb:   cmdb,
	}
}

func ListPage(page *Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err := page.All(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		root := make(map[string]any)
		root["Devices"] = devices

		page.tpl.ExecuteTemplate(w, "list.tpl", root)
	}
}

func AdminPage(page *Page) http.Handler {
	admin := chi.NewRouter()
	admin.Get("/", ListPage(page))
	return admin
}
