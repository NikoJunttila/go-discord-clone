package core

import (
	"discord/internal/handlers/json"
	logger "discord/pkg/logging"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (c *CoreHandler) GetHTML(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	_, err := c.Queries.CreateFoo(r.Context(), title)
	if err != nil {
		slog.Warn("db err", "error", err)
	}
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := struct {
		Title string
		Items []string
	}{
		Title: fmt.Sprintf("<h1>%s</h1>", title),
		Items: []string{"xdd", "item2", "cs"},
	}
	tmpl.Execute(w, data)
}

var Err = errors.New("test error")

func (c *CoreHandler) Example(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Error(ctx, "test log", Err)
	json.RespondWithJSON(ctx, w, 200, "xdd")
}
