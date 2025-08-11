package app

import (
	"context"
	"discord/internal/db"
	logger "discord/pkg/logging"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Connection *pgxpool.Pool
	Queries    *db.Queries
	ctx        context.Context
	cancel     context.CancelFunc
	Logger     *slog.Logger
	// more shared deps (e.g. Redis, Mailer, Config)
}

func Setup(parentCtx context.Context) *App {
	// Create app-specific context that can be cancelled independently
	ctx, cancel := context.WithCancel(parentCtx)
	log := logger.Setup()

	// Add timeout for database initialization
	initCtx, initCancel := context.WithTimeout(ctx, 30*time.Second)
	defer initCancel()

	dbConn, err := db.NewDatabase(initCtx)
	if err != nil {
		cancel() // Clean up our context
		logger.Fatal(ctx, "error", err)
		return nil
	}

	// Test the connection with context
	if err := dbConn.Ping(initCtx); err != nil {
		dbConn.Close()
		cancel()
		logger.Fatal(ctx, "error", err)
		return nil
	}

	return &App{
		Connection: dbConn,
		Queries:    db.New(dbConn),
		Logger:     log,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Context returns the app's context for background operations
func (a *App) Context() context.Context {
	return a.ctx
}

// Cleanup releases resources and cancels background operations
func (a *App) Cleanup() {
	slog.Info("shutting down app resources...")

	// Cancel context first to signal background operations to stop
	if a.cancel != nil {
		a.cancel()
	}

	// Close database connection
	if a.Connection != nil {
		a.Connection.Close()
		slog.Info("database connection closed")
	}
	slog.Info("Clean up complete")
}

// Example method showing how to use the app context for background work
func (a *App) StartBackgroundJobs() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Do some periodic work with database
				a.performPeriodicCleanup()
			case <-a.ctx.Done():
				slog.Info("background jobs stopping due to context cancellation")
				return
			}
		}
	}()
}

func (a *App) performPeriodicCleanup() {
	// Example: cleanup old records with context
	ctx, cancel := context.WithTimeout(a.ctx, 30*time.Second)
	defer cancel()

	// Use a.Queries with context for database operations
	_ = ctx // Use this context for your database queries
	slog.Info("performing periodic cleanup")
}
