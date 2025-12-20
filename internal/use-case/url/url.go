package url_usecase

import (
	"context"
	"errors"
	"main/internal/domain/url"
	url_dto "main/internal/dto/url"
	repo "main/internal/repository/url"
	"time"
)

type URLUseCase struct {
	Repository *repo.URLRepository
}

func New(repository *repo.URLRepository) *URLUseCase {
	return &URLUseCase{
		Repository: repository,
	}
}

func (u *URLUseCase) Create(dto url_dto.CreateURLDTO) error {
	expiredAt := time.Now().Add(48 * time.Hour)

	url, err := url.New(dto.Slug, dto.OriginalURL, &expiredAt)
	if err != nil {
		return err
	}

	err = u.Repository.Create(context.Background(), url)
	if err != nil {
		return err
	}

	return nil
}

func (u *URLUseCase) Redirect(dto url_dto.RedirectURLDTO) (*url.URL, error) {
	url, err := u.Repository.GetBySlug(context.Background(), dto.Slug)
	if err != nil {
		return nil, err
	}

	if expired := url.IsExpired(); expired {
		return nil, errors.New("URL has expired")
	}

	return url, nil
}
