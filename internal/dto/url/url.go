package url_dto

type CreateURLDTO struct {
	OriginalURL string `json:"original_url"`
	Slug        string `json:"slug"`
}
