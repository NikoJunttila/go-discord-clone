package core

import (
	"github.com/go-chi/chi/v5"
)

func (c *CoreHandler) Prefix() string {
	return "/api"
}
func (c *CoreHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/example", c.GetExample)
	r.Get("/ws", c.GetWSChat)
	return r
}
