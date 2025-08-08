package app

import (
	"discord/internal/db"
	"log"
	"context"
)

type App struct {
	Queries *db.Queries
	// Logger  *log.Logger
	// more shared deps (e.g. Redis, Mailer, Config)
}

func Init()*App{
	ctx := context.Background()
	dbConn, err := db.NewDatabase(ctx)
	if err != nil {
		log.Fatal("failed to create db conn, ", err)
		return nil
	}
 return &App{
		Queries: db.New(dbConn),
	}
}
