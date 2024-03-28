package model

type Shortening struct {
	Key string
	URL string
}

type ShortenInput struct {
	RawURL     string
	Identifier string
}
