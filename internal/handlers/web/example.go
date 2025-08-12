package web

import (
	"discord/static"
	"fmt"
	"html/template"
	"net/http"
	// "github.com/go-chi/chi/v5"
)

type pageData struct {
	Title template.HTML
	Items []string
}

func (h *WebHandler) GetExample(w http.ResponseWriter, r *http.Request) {
	// title := chi.URLParam(r, "title")
	title := r.URL.Query().Get("title")
	if title == "" {
		title = "xdd"
	}
	data := pageData{
		Title: template.HTML(fmt.Sprintf("<h1>%s</h1>", title)),
		Items: []string{"xdd", "item2", "cs"},
	}
	static.Templates.ExecuteTemplate(w, "example.html", data)
}
