package main

import (
	"fmt"
	"log"
	"main/internal/config"
	url_handler "main/internal/handler/url"
	"main/internal/metrics"
	url_repository "main/internal/repository/url"
	"main/internal/router"
	url_usecase "main/internal/use-case/url"
	"main/pkg/pgx"
	"main/pkg/redis"
	"main/pkg/trace"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	metrics.Register()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		http.ListenAndServe(fmt.Sprintf(":%d", 2223), nil)
	}()

	database := pgx.New(config.Env.Postgres.URI)
	redis := redis.New(config.Env.Redis.Address, config.Env.Redis.Password, config.Env.Redis.DB)

	repo := url_repository.New(database)
	usecase := url_usecase.New(repo, redis)
	handler := url_handler.New(usecase)
	router := router.New(handler)

	shutdown := trace.New("url-shortener")
	defer shutdown()

	log.Println("HTTP server running on :8080")
	http.ListenAndServe(":8080", router)
}
