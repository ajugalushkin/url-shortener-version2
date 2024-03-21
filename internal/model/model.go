package model

import "errors"

type Shortening struct {
	Key string
	URL string
}

type ShortenInput struct {
	RawURL     string
	Identifier string
}

var (
	ErrNotFound         = errors.New("not found")
	ErrIdentifierExists = errors.New("identifier already exists")
)
