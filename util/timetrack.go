package util

import (
	"time"
	"log/slog"
)

func TimeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    slog.Info("%s took %s", name, elapsed)
}
