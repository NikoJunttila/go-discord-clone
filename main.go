package main

import (
	"context"
	"discord/internal/app"
	"discord/internal/router"
	logger "discord/pkg/logging"
	"discord/util"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, falling back to OS environment variables")
	}
	// Create root context that will be cancelled on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	coreApp := app.Setup(ctx)
	defer coreApp.Cleanup()

	// HTTP server with custom settings
	port := util.GetEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router.SetupRouter(coreApp),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for OS interrupts
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		slog.Info("Server listening...", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "shutting Shutting down gracefully...", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	slog.Info("Shutting down gracefully...")

	// Cancel root context to signal all background operations to stop
	cancel()

	// Create a deadline for HTTP server shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Stop HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal(ctx, "Server forced to shutdown: %v", err)
	}

	slog.Info("Server exited cleanly")
}
