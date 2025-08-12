package core

import (
	"discord/internal/app"
	"discord/internal/db"
	"log/slog"
)

type CoreHandler struct {
	DB  *db.Queries
	log *slog.Logger
}

func NewCoreHandler(app *app.App) *CoreHandler {
	return &CoreHandler{DB: app.Queries, log: app.Logger}
}
