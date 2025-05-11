package page

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/doggy/pf"
)

type Page struct {
	tpl          *template.Template
	deviceLoader pf.DeviceLoader
}

func NewPage(loader pf.DeviceLoader) *Page {
	return &Page{
		tpl:          template.Must(template.ParseFS(tplFS, "tpl/*.tpl")),
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
		root["Devices"] = devices

		page.tpl.ExecuteTemplate(w, "list.tpl", root)
	}
}

func AdminPage(page *Page) http.Handler {
	admin := chi.NewRouter()
	admin.Get("/", ListPage(page))
	return admin
}
