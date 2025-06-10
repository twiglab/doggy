package page

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/doggy/pf"
)

type Loader interface {
	All(context.Context) ([]pf.CameraUpload, error)
}

type Item struct {
	Upload pf.CameraUpload
	TTL    time.Time
}
type Page struct {
	tpl    *template.Template
	Loader Loader
	Touch  pf.Toucher
}

func (v *Page) All(ctx context.Context) ([]Item, error) {
	uploads, err := v.Loader.All(ctx)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, u := range uploads {
		ttl, _ := v.Touch.Last(u.UUID)
		items = append(items, Item{TTL: ttl, Upload: u})
	}

	return items, nil
}

func NewPage(loader Loader, l pf.Toucher) *Page {
	return &Page{
		tpl:    template.Must(template.ParseFS(tplFS, "tpl/*.tpl")),
		Loader: loader,
		Touch:  l,
	}
}

func ListPage(page *Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := page.All(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		root := make(map[string]any)
		root["Items"] = items

		page.tpl.ExecuteTemplate(w, "list.tpl", root)
	}
}

func AdminPage(page *Page) http.Handler {
	admin := chi.NewRouter()
	admin.Get("/", ListPage(page))
	return admin
}
