package model

type Shortening struct {
	Key string
	URL string
}

type File struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ShortenInput struct {
	RawURL     string
	Identifier string
}
