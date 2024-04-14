package storage

import (
	"errors"
	"sync"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
)

type InMemory struct {
	m sync.Map
}

func NewInMemory() *InMemory {
	return &InMemory{}
}

func (s *InMemory) Put(shortening dto.Shortening) (*dto.Shortening, error) {
	if _, exists := s.m.Load(shortening.ShortURL); exists {
		return nil, errors.New("identifier already exists")
	}

	s.m.Store(shortening.ShortURL, shortening)

	return &shortening, nil
}

func (s *InMemory) PutList(list dto.ShorteningList) error {
	for _, shortening := range list {
		_, err := s.Put(shortening)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *InMemory) Get(identifier string) (*dto.Shortening, error) {
	v, ok := s.m.Load(identifier)
	if !ok {
		return nil, errors.New("not found")
	}

	shortening := v.(dto.Shortening)

	return &shortening, nil
}
