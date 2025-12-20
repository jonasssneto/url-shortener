package url_handler

import (
	"encoding/json"
	url_dto "main/internal/dto/url"
	usecase "main/internal/use-case/url"
	"main/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	Usecase *usecase.URLUseCase
	Logger  *logger.Logger
}

func New(usecase *usecase.URLUseCase) *URLHandler {
	return &URLHandler{
		Usecase: usecase,
		Logger:  logger.New("url-handler"),
	}
}

func (u *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto url_dto.CreateURLDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.Logger.Debugf("Received Create request: %+v\n", dto)

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

	u.Logger.Debugf("Received Redirect request for slug: %s", slug)

	url, err := u.Usecase.Redirect(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
