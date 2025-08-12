package web

import (
	"github.com/go-chi/chi/v5"
)

func (w *WebHandler) Prefix() string {
	return "/web"
}

func (w *WebHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/example", w.GetExample)
	r.Get("/chat", w.GetChat)
	return r
}
