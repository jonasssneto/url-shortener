package pgx

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(uri string) *pgxpool.Pool {
	database, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return database
}
