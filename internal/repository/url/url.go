package url_repository

import (
	"context"
	"main/internal/domain/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Conn
}

func New(db *pgxpool.Conn) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (u *URLRepository) Create(ctx context.Context, url *url.URL) error {
	query := `INSERT INTO urls (id, slug, original_url, created_at, expired_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := u.db.Exec(ctx, query,
		url.ID,
		url.Slug,
		url.OriginalURL,
		url.CreatedAt,
		url.ExpiredAt,
		url.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *URLRepository) GetBySlug(ctx context.Context, slug string) (*url.URL, error) {
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
		return nil, err
	}

	return &fetchedURL, nil
}
