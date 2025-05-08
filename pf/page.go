package pf

import (
	"html/template"
	"net/http"

	"embed"

	"github.com/go-chi/chi/v5"
)

//go:embed tpl
var tplFS embed.FS

type Page struct {
	Tpl          *template.Template
	deviceLoader DeviceLoader
}

func NewPage(loader DeviceLoader) *Page {
	return &Page{
		Tpl:          template.Must(template.ParseFS(tplFS, "tpl/*.tpl")),
		deviceLoader: loader,
	}
}

func ListPage(page *Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err := page.deviceLoader.All(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		root := make(map[string]any)
		root["devices"] = devices

		page.Tpl.ExecuteTemplate(w, "list", root)
	}
}

func AdminPage(page *Page) http.Handler {
	admin := chi.NewRouter()
	admin.Get("/", ListPage(page))
	return admin
}
