// Package logger contains all logging logic using slog.
package logger

import (
	"context"
	"discord/util"
	"log/slog"
	"os"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

var defaultLogger *slog.Logger

// withRequestContext adds common context fields to the logger
func withRequestContext(ctx context.Context) *slog.Logger {
	logger := defaultLogger

	// Add request ID if available
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		logger = logger.With("request_id", reqID)
	}
	//Add these later if necessary for now keep commented out
	// // Add user ID if available (common pattern)
	// if userID := getUserID(ctx); userID != "" {
	// 	logger = logger.With("user_id", userID)
	// }
	//
	// // Add trace ID if available (for distributed tracing)
	// if traceID := getTraceID(ctx); traceID != "" {
	// 	logger = logger.With("trace_id", traceID)
	// }

	return logger
}

// Info logs an info message with request_id automatically injected from context
func Info(ctx context.Context, msg string, args ...any) {
	withRequestContext(ctx).Info(msg, args...)
}

// Error logs an error with the request_id (if available in ctx).
func Error(ctx context.Context, msg string, err error) {
	if err == nil {
		return
	}
	withRequestContext(ctx).Error(msg, "error", err)
}

// Warn logs a warning message using the logger from context.
func Warn(ctx context.Context, msg string, args ...any) {
	withRequestContext(ctx).Warn(msg, args...)
}

// Fatal logs a fatal error message using the logger from context and exits the program.
func Fatal(ctx context.Context, msg string, err error, args ...any) {
	withRequestContext(ctx).Error(msg, "error", err)
	os.Exit(1)
}

// Debug logs a debug message using the logger from context.
func Debug(ctx context.Context, msg string, args ...any) {
	withRequestContext(ctx).Debug(msg, args...)
}

// ===== Setup function with Loki integration =====
// It configures logging based on environment variables:
// - LOG_LEVEL: debug, info, warn, error (default: info)
// - LOG_FORMAT: json, text (default: json)
// - SERVICE_NAME: adds service field to all logs
// - LOG_ADD_SOURCE: true/false (default: false in production)
func Setup() *slog.Logger {
	level := slog.LevelInfo
	switch strings.ToLower(util.GetEnv("LOG_LEVEL", "info")) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	addSource := strings.ToLower(util.GetEnv("LOG_ADD_SOURCE", "false")) == "true"

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	}

	serviceName := util.GetEnv("SERVICE_NAME", "go_project")
	lokiURL := util.GetEnv("LOKI_URL", "")

	// Always log to stdout

	if lokiURL != "" {
		// Log to both stdout and Loki
		stdoutHandler := slog.NewJSONHandler(os.Stdout, opts)
		handler := newMultiHandler(stdoutHandler, newLokiHandler(lokiURL, serviceName))
		defaultLogger = slog.New(handler).With("service", serviceName)
	} else {
		// Log only to stdout
		if util.GetEnv("LOG_FORMAT", "json") == "text" {
			stdoutHandler := slog.NewTextHandler(os.Stdout, opts)
			defaultLogger = slog.New(stdoutHandler)
		} else {
			stdoutHandler := slog.NewJSONHandler(os.Stdout, opts)
			defaultLogger = slog.New(stdoutHandler)
		}
		// defaultLogger = slog.New(stdoutHandler).With("service", serviceName)
	}

	slog.SetDefault(defaultLogger)

	defaultLogger.Info("global logger initialized",
		"level", level.String(),
		"add_source", addSource,
		"loki_enabled", lokiURL != "")

	return defaultLogger
}
