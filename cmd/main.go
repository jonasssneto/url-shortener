package main

import (
	"context"
	"log"
	"main/internal/config"
	url_handler "main/internal/handler/url"
	url_repository "main/internal/repository/url"
	"main/internal/router"
	url_usecase "main/internal/use-case/url"
	"net/http"

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

	repo := url_repository.New(conn)
	usecase := url_usecase.New(repo)
	handler := url_handler.New(usecase)
	router := router.New(handler)

	log.Println("HTTP server running on :8080")
	http.ListenAndServe(":8080", router)
}
