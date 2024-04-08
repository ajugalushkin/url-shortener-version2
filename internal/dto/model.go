package dto

type Shortening struct {
	Key string
	URL string
}

type ShortenInput struct {
	RawURL     string
	Identifier string
}
