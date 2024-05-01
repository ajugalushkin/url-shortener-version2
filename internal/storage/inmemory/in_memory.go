package inmemory

import (
	"context"
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

func (r *InMemory) Put(ctx context.Context, shortening dto.Shortening) (*dto.Shortening, error) {
	if _, exists := r.m.Load(shortening.ShortURL); exists {
		return nil, errors.New("identifier already exists")
	}

	r.m.Store(shortening.ShortURL, shortening)

	return &shortening, nil
}

func (r *InMemory) PutList(ctx context.Context, list dto.ShorteningList) error {
	for _, shortening := range list {
		_, err := r.Put(ctx, shortening)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *InMemory) Get(ctx context.Context, identifier string) (*dto.Shortening, error) {
	v, ok := r.m.Load(identifier)
	if !ok {
		return nil, errors.New("not found")
	}

	shortening := v.(dto.Shortening)

	return &shortening, nil
}

func (r *InMemory) GetListByUser(ctx context.Context, userID string) (*dto.ShorteningList, error) {
	return &dto.ShorteningList{}, nil
}
