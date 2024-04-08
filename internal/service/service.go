package service

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/shorten"
	"github.com/google/uuid"
)

type PutGetter interface {
	Put(urlData dto.Shortening) (*dto.Shortening, error)
	Get(id string) (*dto.Shortening, error)
}

type Service struct {
	storage PutGetter
}

func NewService(storage PutGetter) *Service {
	return &Service{storage: storage}
}

func (s *Service) Shorten(input dto.ShortenInput) (*dto.Shortening, error) {
	var (
		id         = uuid.New().ID()
		identifier = input.Identifier
	)
	if identifier == "" {
		identifier = shorten.Shorten(id)
	}

	newShortening := dto.Shortening{
		Key: identifier,
		URL: input.RawURL,
	}

	shortening, err := s.storage.Get(newShortening.Key)
	if err != nil {
		shortening, err = s.storage.Put(newShortening)
		if err != nil {
			return nil, err
		}
	}

	return shortening, nil
}

func (s *Service) Redirect(identifier string) (string, error) {
	shortening, err := s.storage.Get(identifier)
	if err != nil {
		return "", err
	}

	return shortening.URL, nil
}
