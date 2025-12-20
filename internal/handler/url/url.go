package url_handler

import (
	"encoding/json"
	"log"
	url_dto "main/internal/dto/url"
	usecase "main/internal/use-case/url"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	Usecase *usecase.URLUseCase
}

func New(usecase *usecase.URLUseCase) *URLHandler {
	return &URLHandler{
		Usecase: usecase,
	}
}

func (u *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto url_dto.CreateURLDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received Create URL request: %+v\n", dto)

	err = u.Usecase.Create(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	dto := url_dto.RedirectURLDTO{
		Slug: slug,
	}

	log.Printf("Received Redirect request for slug: %s\n", slug)

	url, err := u.Usecase.Redirect(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
