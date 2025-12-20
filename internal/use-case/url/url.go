package url_usecase

import (
	"context"
	"errors"
	"main/internal/domain/url"
	url_dto "main/internal/dto/url"
	repo "main/internal/repository/url"
	"main/pkg/logger"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type URLUseCase struct {
	Repository *repo.URLRepository
	Logger     *logger.Logger
}

func New(repository *repo.URLRepository) *URLUseCase {
	return &URLUseCase{
		Repository: repository,
		Logger:     logger.New("url-usecase"),
	}
}

func (u *URLUseCase) Create(ctx context.Context, dto url_dto.CreateURLDTO) error {
	tracer := otel.Tracer("url-usecase")
	ctx, span := tracer.Start(ctx, "create")
	defer span.End()

	expiredAt := time.Now().Add(48 * time.Hour)
	u.Logger.Debugf("Creating URL with slug: %s, original URL: %s, expires at: %s", dto.Slug, dto.OriginalURL, expiredAt)

	url, err := url.New(dto.Slug, dto.OriginalURL, &expiredAt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid url data")
		return err
	}

	err = u.Repository.Create(ctx, url)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "repository create failed")
		return err
	}

	return nil
}

func (u *URLUseCase) Redirect(ctx context.Context, dto url_dto.RedirectURLDTO) (*url.URL, error) {
	tracer := otel.Tracer("url-usecase")
	ctx, span := tracer.Start(ctx, "redirect")
	defer span.End()

	u.Logger.Debugf("Redirecting for slug: %s", dto.Slug)
	url, err := u.Repository.GetBySlug(ctx, dto.Slug)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "repository get failed")
		return nil, err
	}

	if expired := url.IsExpired(); expired {
		u.Logger.Debugf("URL with slug %s has expired", dto.Slug)
		span.SetStatus(codes.Error, "url expired")
		return nil, errors.New("URL has expired")
	}

	return url, nil
}
