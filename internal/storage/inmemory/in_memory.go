package inmemory

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
)

// InMemory структура репозитория.
type InMemory struct {
	m sync.Map
}

// NewInMemory конструктор.
func NewInMemory() *InMemory {
	return &InMemory{}
}

// Put метод сохраняет данные в память.
func (r *InMemory) Put(ctx context.Context, shortening dto.Shortening) (*dto.Shortening, error) {
	if _, exists := r.m.Load(shortening.ShortURL); exists {
		err := errors.New("identifier already exists")
		logger.LogFromContext(ctx).Debug("InMemory.Put Load Error",
			zap.Error(err))
		return nil, err
	}

	r.m.Store(shortening.ShortURL, shortening)

	logger.LogFromContext(ctx).Debug("InMemory.Put Store Success")
	return &shortening, nil
}

// PutList метод сохраняет список данных в память.
func (r *InMemory) PutList(ctx context.Context, list dto.ShorteningList) error {
	for _, shortening := range list {
		_, err := r.Put(ctx, shortening)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get метод позволяет получить данные по короткому URL
func (r *InMemory) Get(ctx context.Context, identifier string) (*dto.Shortening, error) {
	v, ok := r.m.Load(identifier)
	if !ok {
		return nil, errors.New("not found")
	}

	shortening := v.(dto.Shortening)

	return &shortening, nil
}

// GetListByUser метод позволяет получить список URL для конкретного пользователя
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

// DeleteUserURL метод позволяет удалить список URL для конкретного пользователя
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
