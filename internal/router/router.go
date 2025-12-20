package router

import (
	url_handler "main/internal/handler/url"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func New(urlHandler *url_handler.URLHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/url", urlHandler.Create)
	r.Get("/{slug}", urlHandler.Redirect)

	return r
}
