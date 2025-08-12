package web

import (
	"discord/static"
	"html/template"
	"net/http"
)

type ChatPageData struct {
	Title template.HTML
}

func (h *WebHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	data := ChatPageData{
		Title: template.HTML("WebSocket Chat"),
	}
	if err := static.Templates.ExecuteTemplate(w, "chat.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
