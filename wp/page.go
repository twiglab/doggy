package wp

import (
	"html/template"
	"net/http"
)

type Page struct {
	Tpl *template.Template
}

func NewPage() *Page {
	return &Page{
		Tpl: template.Must(template.ParseFS(tplFS, "tpl/*.tpl")),
	}
}
func ListPage(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
