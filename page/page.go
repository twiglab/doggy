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
	GetAll(context.Context) ([]pf.Channel, error)
}

type Item struct {
	Upload pf.Channel
	TTL    time.Time
}
type Page struct {
	tpl    *template.Template
	Loader Loader
	Touch  pf.Toucher
}

func (v *Page) All(ctx context.Context) ([]Item, error) {
	uploads, err := v.Loader.GetAll(ctx)
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
	t := template.Must(template.New("page").Funcs(build()).ParseFS(tplFS, "tpl/*.tpl"))
	return &Page{
		tpl:    t,
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

func datetime(t time.Time) string {
	return t.Format(time.DateTime)
}

func alarm(t time.Time) string {
	d := time.Until(t).Abs()
	if d < 90*time.Second {
		return "ok"
	}
	if d > 3*time.Minute {
		return "alarm"
	}

	return "warn"
}

func build() template.FuncMap {
	return template.FuncMap{
		"alarm":    alarm,
		"datetime": datetime,
	}
}
