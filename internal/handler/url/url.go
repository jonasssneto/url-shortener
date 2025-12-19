package url_handler

import (
	"encoding/json"
	url_dto "main/internal/dto/url"
	usecase "main/internal/use-case/url"
	"net/http"
)

type URLHandler struct {
	Usecase usecase.URLUseCase
}

func New(usecase usecase.URLUseCase) *URLHandler {
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

	err = u.Usecase.Create(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
