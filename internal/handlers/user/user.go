package user

import (
	"discord/internal/app"
	"discord/internal/db"
	"net/http"
	"log/slog"
	"github.com/go-chi/chi/v5"
	"html/template"
	"fmt"
)

type UserHandler struct {
	Queries *db.Queries
}

type pageData struct {
	Title string
	Items []string
}
func NewUserHandler(app *app.App) *UserHandler{
	return &UserHandler{Queries: app.Queries}
}

func (h *UserHandler)GetBasic(w http.ResponseWriter, r *http.Request){
	title := chi.URLParam(r, "title")
	_, err := h.Queries.CreateFoo(r.Context(),title)
  if err != nil {
	slog.Warn("db err", "error", err)
  }
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := pageData{
		Title: fmt.Sprintf("<h1>%s</h1>",title),
		Items: []string{"xdd", "item2", "cs"},
	}
	tmpl.Execute(w, data)
}
