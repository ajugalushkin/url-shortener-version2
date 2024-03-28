package model

//go:generate easyjson -all json.go
type (
	Shorten struct {
		URL string `json:"url"`
	}
	ShortenResult struct {
		Result string `json:"result"`
	}
	File struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
)
