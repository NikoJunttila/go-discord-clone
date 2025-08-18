package core

import (
	"discord/internal/app"
	"discord/internal/service/user"
	"log/slog"
)

type CoreHandler struct {
	Log         *slog.Logger
	Chat        *ChatHub
	UserService user.Service
}

func NewCoreHandler(app *app.App) *CoreHandler {
	hub := NewChatHub()
	go hub.Run(app.Context())
	return &CoreHandler{
		Log:         app.Logger,
		Chat:        hub,
		UserService: app.UserService,
	}
}
