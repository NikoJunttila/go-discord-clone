package util

import (
	"log/slog"
	"os"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Error("Failed to find env", "env", key)
		return defaultValue
	}
	return value
}

