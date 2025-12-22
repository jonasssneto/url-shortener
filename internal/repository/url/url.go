package url_repository

import (
	"context"
	"main/internal/domain/url"
	"main/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type URLRepository struct {
	db     *pgxpool.Pool
	Logger *logger.Logger
}

func New(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		db:     db,
		Logger: logger.New("url-repository"),
	}
}

func (u *URLRepository) Create(ctx context.Context, url *url.URL) error {
	tracer := otel.Tracer("url-repository")
	ctx, span := tracer.Start(ctx, "create-db")
	defer span.End()

	span.SetAttributes(
		attribute.String("url.slug", url.Slug),
		attribute.String("url.original", url.OriginalURL),
		attribute.String("url.expired_at", url.ExpiredAt.Format(time.RFC3339)),
	)

	u.Logger.Debugf("Inserting URL into database: %+v", url)
	query := `INSERT INTO urls (id, slug, original_url, expired_at)
		  VALUES ($1, $2, $3, $4)`

	_, err := u.db.Exec(ctx, query,
		url.ID,
		url.Slug,
		url.OriginalURL,
		url.ExpiredAt,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "db exec failed")
		return err
	}

	span.SetAttributes(
		attribute.Bool("url.created", true),
	)

	return nil
}

func (u *URLRepository) GetBySlug(ctx context.Context, slug string) (*url.URL, error) {
	tracer := otel.Tracer("url-repository")
	ctx, span := tracer.Start(ctx, "get-by-slug-db")
	defer span.End()

	span.SetAttributes(
		attribute.String("url.slug", slug),
	)
	query := `SELECT id, original_url,expired_at
			  FROM urls
			  WHERE slug = $1`

	row := u.db.QueryRow(ctx, query, slug)

	var fetchedURL url.URL
	err := row.Scan(
		&fetchedURL.ID,
		&fetchedURL.OriginalURL,
		&fetchedURL.ExpiredAt,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "db query failed")
		return nil, err
	}

	span.SetAttributes(
		attribute.String("url.original", fetchedURL.OriginalURL),
		attribute.String("url.expired_at", fetchedURL.ExpiredAt.Format(time.RFC3339)),
	)

	urlDomain := &url.URL{
		ID:          fetchedURL.ID,
		Slug:        slug,
		OriginalURL: fetchedURL.OriginalURL,
		ExpiredAt:   fetchedURL.ExpiredAt,
	}

	span.SetAttributes(
		attribute.Bool("url.found", true),
	)

	return urlDomain, nil
}
