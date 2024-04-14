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

func (s *Service) Shorten(input dto.Shortening) (*dto.Shortening, error) {
	var (
		id         = uuid.New().ID()
		identifier = input.ShortURL
	)
	if identifier == "" {
		identifier = shorten.Shorten(id)
	}

	newShortening := dto.Shortening{
		ShortURL:      identifier,
		OriginalURL:   input.OriginalURL,
		CorrelationId: input.CorrelationId,
	}

	shortening, err := s.storage.Get(newShortening.ShortURL)
	if err != nil {
		shortening, err = s.storage.Put(newShortening)
		if err != nil {
			return nil, err
		}
	}

	return shortening, nil
}

func (s *Service) ShortenList(input dto.ShortenListInput) (*dto.ShorteningList, error) {
	var resultList dto.ShorteningList
	for _, item := range input {
		newShortening := dto.Shortening{
			ShortURL:      shorten.Shorten(uuid.New().ID()),
			OriginalURL:   item.OriginalURL,
			CorrelationId: item.CorrelationId,
		}

		resultShortening, err := s.storage.Get(newShortening.ShortURL)
		if err != nil {
			resultShortening, err = s.storage.Put(newShortening)
			if err != nil {
				return nil, err
			}
		}

		resultList = append(resultList, *resultShortening)
	}
	return &resultList, nil
}

func (s *Service) Redirect(identifier string) (string, error) {
	shortening, err := s.storage.Get(identifier)
	if err != nil {
		return "", err
	}

	return shortening.OriginalURL, nil
}
