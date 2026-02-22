package url_usecase

import (
	"context"
	"errors"
	"main/internal/domain/url"
	url_dto "main/internal/dto/url"
	repo "main/internal/repository/url"
	"main/pkg/logger"
	"main/pkg/redis"
	"time"

	redis_client "github.com/go-redis/redis"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type URLUseCase struct {
	Repository *repo.URLRepository
	Cache      *redis.Client
	Logger     *logger.Logger
}

func New(repository *repo.URLRepository, cache *redis.Client) *URLUseCase {
	return &URLUseCase{
		Repository: repository,
		Cache:      cache,
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

	span.SetAttributes(
		attribute.String("slug", dto.Slug),
	)

	cachedURL, err := u.Cache.Conn.Get(dto.Slug).Result()
	if err == nil {
		span.SetAttributes(attribute.Bool("cache.hit", true))

		u.Logger.Debugf("Cache hit for slug: %s", dto.Slug)

		span.SetStatus(codes.Ok, "")
		return &url.URL{
			Slug:        dto.Slug,
			OriginalURL: cachedURL,
		}, nil
	}

	if err != redis_client.Nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("cache.error", true))
	}

	span.SetAttributes(attribute.Bool("cache.hit", false))

	urlModel, err := u.Repository.GetBySlug(ctx, dto.Slug)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "repository failure")
		return nil, err
	}

	if urlModel.IsExpired() {
		span.SetAttributes(attribute.Bool("url.expired", true))

		u.Logger.Debugf("URL with slug %s has expired", dto.Slug)

		return nil, errors.New("URL has expired")
	}

	if err := u.Cache.Conn.Set(dto.Slug, urlModel.OriginalURL, 15*time.Minute).Err(); err != nil {
		u.Logger.Errorf("Failed to cache URL for slug %s: %v", dto.Slug, err)
		span.RecordError(err)
	}

	span.SetStatus(codes.Ok, "")
	return urlModel, nil
}
