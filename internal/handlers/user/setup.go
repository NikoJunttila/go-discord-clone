package user

import (
	"discord/internal/app"
	"discord/internal/db"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	Queries *db.Queries
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{Queries: app.Queries}
}

func (h *UserHandler) Prefix() string {
	return "/api"
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/title/{title}", h.GetBasic)
	return r
}
