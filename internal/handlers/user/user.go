package user

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
)

type pageData struct {
	Title string
	Items []string
}

func (h *UserHandler) GetBasic(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	_, err := h.Queries.CreateFoo(r.Context(), title)
	if err != nil {
		slog.Warn("db err", "error", err)
	}
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := pageData{
		Title: fmt.Sprintf("<h1>%s</h1>", title),
		Items: []string{"xdd", "item2", "cs"},
	}
	tmpl.Execute(w, data)
}
