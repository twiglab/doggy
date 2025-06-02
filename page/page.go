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

type Page struct {
	tpl    *template.Template
	loader Loader
}

func NewPage(loader Loader) *Page {
	return &Page{
		tpl:    template.Must(template.ParseFS(tplFS, "tpl/*.tpl")),
		loader: loader,
	}
}

func ListPage(page *Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err := page.loader.All(r.Context())
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
