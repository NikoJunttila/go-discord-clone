package core

import (
	"discord/internal/app"
	"discord/internal/db"
	"log/slog"
)

type CoreHandler struct {
	DB   *db.Queries
	Log  *slog.Logger
	Chat *ChatHub
}

func NewCoreHandler(app *app.App) *CoreHandler {
	hub := NewChatHub()
	go hub.Run(app.Context())
	return &CoreHandler{DB: app.Queries, Log: app.Logger, Chat: hub}
}
