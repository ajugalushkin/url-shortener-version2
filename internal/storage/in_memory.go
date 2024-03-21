package storage

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"sync"
)

type InMemory struct {
	m sync.Map
}

func NewInMemory() *InMemory {
	return &InMemory{}
}

func (s *InMemory) Put(shortening model.Shortening) (*model.Shortening, error) {
	if _, exists := s.m.Load(shortening.Key); exists {
		return nil, model.ErrIdentifierExists
	}

	s.m.Store(shortening.Key, shortening)

	return &shortening, nil
}

func (s *InMemory) Get(identifier string) (*model.Shortening, error) {
	v, ok := s.m.Load(identifier)
	if !ok {
		return nil, model.ErrNotFound
	}

	shortening := v.(model.Shortening)

	return &shortening, nil
}
