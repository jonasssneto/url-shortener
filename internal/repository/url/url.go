package url_repository

import (
	"context"
	"log"
	"main/internal/domain/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (u *URLRepository) Create(ctx context.Context, url *url.URL) error {
	log.Println("Inserting URL into database:", url)
	query := `INSERT INTO urls (id, slug, original_url, expired_at)
          VALUES ($1, $2, $3, $4)`

	_, err := u.db.Exec(ctx, query,
		url.ID,
		url.Slug,
		url.OriginalURL,
		url.ExpiredAt,
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

	urlDomain := &url.URL{
		ID:          fetchedURL.ID,
		Slug:        slug,
		OriginalURL: fetchedURL.OriginalURL,
		ExpiredAt:   fetchedURL.ExpiredAt,
	}

	return urlDomain, nil
}
