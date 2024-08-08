package dto

//go:generate easyjson -all json.go

// User структура cookie
type User struct {
	ID int
}

//easyjson:json
type Shortening struct {
	CorrelationID string `json:"correlation_id" db:"correlation_id" `
	ShortURL      string `json:"short_url" db:"short_url"`
	OriginalURL   string `json:"original_url" db:"original_url"`
	UserID        string `json:"user_id" db:"user_id"`
	IsDeleted     bool   `json:"is_deleted" db:"is_deleted"`
}

//easyjson:json
type ShorteningList []struct {
	CorrelationID string `json:"correlation_id" db:"correlation_id" `
	ShortURL      string `json:"short_url" db:"short_url"`
	OriginalURL   string `json:"original_url" db:"original_url"`
	UserID        string `json:"user_id" db:"user_id"`
	IsDeleted     bool   `json:"is_deleted" db:"is_deleted"`
}

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

//easyjson:json
type UserURLListLine struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

//easyjson:json
type UserURLList []struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

//easyjson:json
type URLs []string

type Stats struct {
	URLS  int `json:"URLS" db:"urls"`
	Users int `json:"Users" db:"users"`
}
