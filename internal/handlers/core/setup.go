package core

import (
	"discord/internal/app"
	"discord/internal/db"
	"log/slog"

	"github.com/go-chi/chi/v5"
)

type CoreHandler struct {
	Queries *db.Queries
	log     *slog.Logger
}

func NewCoreHandler(app *app.App) *CoreHandler {
	return &CoreHandler{Queries: app.Queries, log: app.Logger}
}
func (h *CoreHandler) Prefix() string {
	return "/web"
}

func (h *CoreHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/title/{title}", h.GetHTML)
	r.Get("/", h.Example)
	return r
}
