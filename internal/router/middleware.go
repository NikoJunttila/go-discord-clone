package router

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

const slowRequestThreshold = 800 * time.Millisecond

// contextKey is a private type to avoid collisions in context keys.
type contextKey string

const loggerKey contextKey = "logger"

// NewContext returns a new context with the given slog logger.
func NewContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

// FromContext extracts a slog.Logger from the context, or returns the default logger if none is found.
func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return log
	}
	return slog.Default()
}

func RequestLogger(baseLogger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			requestID := middleware.GetReqID(r.Context())

			reqLogger := baseLogger.With(
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote", r.RemoteAddr),
			)

			// Inject the contextual logger into the request
			ctx := NewContext(r.Context(), reqLogger)
			r = r.WithContext(ctx)

			// Serve the request
			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			// Prepare log entry with response details
			entry := reqLogger.With(
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Duration("latency", duration),
			)

			if duration > slowRequestThreshold {
				entry.Warn("Slow request detected")
			}

			switch {
			case ww.Status() >= 500:
				entry.Error("Server error")
			case ww.Status() >= 400:
				entry.Warn("Client error")
			default:
				entry.Info("Request handled successfully")
			}
		})
	}
}
