package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// ===== Loki minimal structs =====

type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type lokiHandler struct {
	lokiURL    string
	service    string
	httpClient *http.Client
}

func newLokiHandler(lokiURL, service string) slog.Handler {
	return &lokiHandler{
		lokiURL:    strings.TrimRight(lokiURL, "/"),
		service:    service,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (h *lokiHandler) Enabled(ctx context.Context, level slog.Level) bool { return true }

func (h *lokiHandler) Handle(ctx context.Context, record slog.Record) error {
	labels := map[string]string{
		"service": h.service,
		"level":   record.Level.String(),
	}
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		labels["request_id"] = reqID
	}

	logEntry := map[string]interface{}{
		"timestamp": record.Time.Format(time.RFC3339Nano),
		"level":     record.Level.String(),
		"message":   record.Message,
		"service":   h.service,
	}
	record.Attrs(func(attr slog.Attr) bool {
		logEntry[attr.Key] = attr.Value.Any()
		return true
	})

	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	pushReq := lokiPushRequest{
		Streams: []lokiStream{
			{
				Stream: labels,
				Values: [][]string{
					{
						strconv.FormatInt(record.Time.UnixNano(), 10),
						string(logJSON),
					},
				},
			},
		},
	}

	go h.sendToLoki(pushReq)
	return nil
}

func (h *lokiHandler) sendToLoki(pushReq lokiPushRequest) {
	body, err := json.Marshal(pushReq)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", h.lokiURL+"/loki/api/v1/push", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	h.httpClient.Do(req) // ignore errors for minimal version
}

func (h *lokiHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *lokiHandler) WithGroup(name string) slog.Handler       { return h }

// ===== MultiHandler =====

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, h := range m.handlers {
		_ = h.Handle(ctx, record) // ignore individual handler errors
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}

func newMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{handlers: handlers}
}
