package storage

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"sync"
)

type inMemory struct {
	m sync.Map
}

func NewInMemory() *inMemory {
	return &inMemory{}
}

func (s *inMemory) Put(shortening model.Shortening) (*model.Shortening, error) {
	if _, exists := s.m.Load(shortening.Key); exists {
		return nil, model.ErrIdentifierExists
	}

	s.m.Store(shortening.Key, shortening)

	return &shortening, nil
}

func (s *inMemory) Get(identifier string) (*model.Shortening, error) {
	v, ok := s.m.Load(identifier)
	if !ok {
		return nil, model.ErrNotFound
	}

	shortening := v.(model.Shortening)

	return &shortening, nil
}
