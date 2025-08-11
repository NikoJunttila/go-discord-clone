// Package logger contains all logging logic using slog.
package logger

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"strings"
)

var defaultLogger *slog.Logger

// withRequestContext adds common context fields to the logger
func withRequestContext(ctx context.Context) *slog.Logger {
	logger := defaultLogger

	// Add request ID if available
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		logger = logger.With("request_id", reqID)
	}

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

// Setup initializes the global slog logger for the application.
// It configures logging based on environment variables:
// - LOG_LEVEL: debug, info, warn, error (default: info)
// - LOG_FORMAT: json, text (default: json)
// - SERVICE_NAME: adds service field to all logs
// - LOG_ADD_SOURCE: true/false (default: false in production)
func Setup() *slog.Logger {
	// Get log level from environment (default: info)
	level := slog.LevelInfo
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		switch strings.ToLower(logLevel) {
		case "debug":
			level = slog.LevelDebug
		case "warn", "warning":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}

	// Determine if source should be added (file:line info)
	addSource := false
	if sourceEnv := os.Getenv("LOG_ADD_SOURCE"); sourceEnv != "" {
		addSource = strings.ToLower(sourceEnv) == "true"
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	}

	// Choose handler based on LOG_FORMAT environment variable
	var handler slog.Handler
	if getLogFormat() == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		// Default to JSON for production
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	defaultLogger = logger

	// Add service name if provided
	if serviceName := os.Getenv("SERVICE_NAME"); serviceName != "" {
		logger = logger.With("service", serviceName)
	}

	// Set as default logger
	slog.SetDefault(logger)

	// Log initialization
	logger.Info("global logger initialized",
		"level", level.String(),
		"format", getLogFormat(),
		"add_source", addSource)
	return logger
}

// getLogFormat returns the current log format for logging purposes
func getLogFormat() string {
	if os.Getenv("LOG_FORMAT") == "text" {
		return "text"
	}
	return "json"
}
