package web

import (
	"discord/internal/app"
	"discord/internal/db"
)

type WebHandler struct {
	DB *db.Queries
}

func NewWebHandler(app *app.App) *WebHandler {
	return &WebHandler{DB: app.Queries}
}
