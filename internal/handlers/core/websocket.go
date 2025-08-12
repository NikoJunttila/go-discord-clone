package core

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow any origin
}

// ChatHub manages clients and broadcasts
type ChatHub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (h *ChatHub) Run(ctx context.Context) {
	for {
		select {
		case msg := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					slog.Error("Web socket write error", "err", err)
					conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (h *ChatHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WS upgrade error", "err", err)
		return
	}

	// Register new client
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	slog.Debug("new client connected", "total", len(h.clients))

	// Listen for messages from this client
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Debug("read error", "err", err)
			break
		}
		if mt == websocket.TextMessage {
			slog.Debug("WS message", "message", msg)
			h.broadcast <- msg
		}
	}

	// Remove client on disconnect
	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()
	conn.Close()
}

func (h *CoreHandler) GetWSChat(w http.ResponseWriter, r *http.Request) {
	h.Chat.HandleWebSocket(w, r)
}
