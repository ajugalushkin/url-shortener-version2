package model

//go:generate easyjson -all json.go
type (
	Shorten struct {
		URL string `json:"url"`
	}
	ShortenResult struct {
		Result string `json:"result"`
	}
)
