package dto

//go:generate easyjson -all json.go

//easyjson:json
type Shortening struct {
	CorrelationID string `json:"correlation_id" db:"correlation_id" `
	ShortURL      string `json:"short_url" db:"short_url"`
	OriginalURL   string `json:"original_url" db:"original_url"`
}

//easyjson:json
type ShorteningList []Shortening

//easyjson:json
type ShortenInput struct {
	URL string `json:"url"`
}

//easyjson:json
type ShortenOutput struct {
	Result string `json:"result"`
}

//easyjson:json
type ShortenListInput []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

//easyjson:json
type ShortenListOutputLine struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

//easyjson:json
type ShortenListOutput []ShortenListOutputLine
