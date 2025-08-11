package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"discord/internal/app"
	"discord/internal/handlers/core"
	"discord/internal/handlers/user"
)

type Routable interface {
	Prefix() string     // e.g. "/web" or "/api"
	Routes() chi.Router // returns a router with all its endpoints
}

func SetupRouter(a *app.App) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(RequestLogger(a.Logger))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Main API version prefix
	v1 := chi.NewRouter()
	v1modules := []Routable{
		core.NewCoreHandler(a),
		user.NewUserHandler(a),
	}
	for _, m := range v1modules {
		v1.Mount(m.Prefix(), m.Routes())
	}
	// Attach v1 to main router
	r.Mount("/v1", v1)

	return r
}
