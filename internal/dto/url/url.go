package url_dto

type CreateURLDTO struct {
	OriginalURL string `json:"original_url"`
	Slug        string `json:"slug"`
}

type RedirectURLDTO struct {
	Slug string `json:"slug"`
}
