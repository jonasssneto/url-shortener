package main

import (
	"context"
	"fmt"
	"log"
	"main/internal/config"
	url_handler "main/internal/handler/url"
	"main/internal/metrics"
	url_repository "main/internal/repository/url"
	"main/internal/router"
	url_usecase "main/internal/use-case/url"
	"main/pkg/trace"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	metrics.Register()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		http.ListenAndServe(fmt.Sprintf(":%d", 2223), nil)
	}()

	repo := url_repository.New(conn)
	usecase := url_usecase.New(repo)
	handler := url_handler.New(usecase)
	router := router.New(handler)

	shutdown := trace.New("url-shortener")
	defer shutdown()

	log.Println("HTTP server running on :8080")
	http.ListenAndServe(":8080", router)
}
