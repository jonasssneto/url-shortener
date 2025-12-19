package url

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID          uuid.UUID  `json:"id"`
	Slug        string     `json:"slug"`
	OriginalURL string     `json:"original_url"`
	CreatedAt   int64      `json:"created_at"`
	ExpiredAt   *time.Time `json:"expired_at"`
	UpdatedAt   int64      `json:"updated_at"`
}

func New(slug, originalURL string, expiredAt *time.Time) (*URL, error) {
	err := validate(slug, originalURL, expiredAt)
	if err != nil {
		return nil, err
	}

	return &URL{
		ID:          uuid.New(),
		Slug:        slug,
		OriginalURL: originalURL,
		ExpiredAt:   expiredAt,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}, nil
}

func validate(slug, originalURL string, expiredAt *time.Time) error {
	if slug == "" {
		return errors.New("slug cannot be empty")
	}

	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return errors.New("invalid original URL")
	}

	if expiredAt != nil && expiredAt.Before(time.Now()) {
		return errors.New("expiration date must be in the future")
	}

	return nil
}

func (u *URL) IsExpired() bool {
	return time.Now().After(*u.ExpiredAt)
}
