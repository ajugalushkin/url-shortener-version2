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
	list := dto.ShorteningList{}
	r.m.Range(func(k, v interface{}) bool {
		item := v.(dto.Shortening)
		if item.UserID == userID {
			list = append(list, item)
		}
		return true
	})
	return &list, nil
}

func (r *InMemory) DeleteUserURL(ctx context.Context, shortURL []string, userID int) {
	for _, value := range shortURL {
		v, ok := r.m.Load(value)
		if !ok {
			continue
		}
		newShortening := v.(dto.Shortening)
		newShortening.IsDeleted = true

		r.m.CompareAndSwap(value, v.(dto.Shortening), newShortening)
	}
}
