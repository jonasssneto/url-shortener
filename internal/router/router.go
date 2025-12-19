package router

import (
	url_handler "main/internal/handler/url"
	"net/http"
)

func New(urlHandler *url_handler.URLHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/url", urlHandler.Create)

	return mux
}
