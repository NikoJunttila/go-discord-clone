package db

import (
	"context"
	"fmt"
	"log"

	migration "discord/sql/schema"
	"discord/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	env := util.GetEnv("ENVIRONMENT", "dev")

	var connStr string
	if env != "prod" {
		dbHost := util.GetEnv("DB_HOST", "localhost")
		dbPort := util.GetEnv("DB_PORT", "5433")
		dbUser := util.GetEnv("DB_USER", "postgres")
		dbPassword := util.GetEnv("DB_PASSWORD", "postgres")
		dbName := util.GetEnv("DB_NAME", "go_chat_db")

		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName,
		)
		fmt.Println("connectiong ",connStr)
	} else {
		connStr = util.GetEnv("CONNECTION_STRING", "")
		if connStr == "" {
			log.Fatal("CONNECTION_STRING must be set in production")
		}
	}
	err := migration.GooseUp(connStr)
	if err != nil {
  return nil, err
  }
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return pool, nil
}

