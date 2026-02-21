package router

import (
	"main/internal/config"
	url_handler "main/internal/handler/url"
	"main/pkg/logger"
	"net/http"

	internal_middleware "main/internal/middleware"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func New(urlHandler *url_handler.URLHandler) http.Handler {
	r := chi.NewRouter()

	logger := logger.New("router")

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.Env.Cors.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: false,
	}))

	r.Use(httprate.LimitByIP(50, 1))
	r.Use(internal_middleware.Metrics())
	r.Use(internal_middleware.Logger(logger))
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/url", urlHandler.Create)
	r.Get("/r/{slug}", urlHandler.Redirect)

	return r
}
