package main

import (
	"context"
	"main/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	conn, err := pgxpool.New(context.Background(), config.Env.Postgres.URI)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

}
