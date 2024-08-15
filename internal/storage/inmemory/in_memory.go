package inmemory

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	userErr "github.com/ajugalushkin/url-shortener-version2/internal/errors"
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
		logger.GetLogger().Debug("InMemory.Put Load Error", zap.Error(err))
		return nil, err
	}

	urlStore, err := r.GetByURL(ctx, shortening.OriginalURL)
	if err != nil {
		return nil, err
	}
	if urlStore.OriginalURL == shortening.OriginalURL {
		return &shortening, errors.Wrapf(userErr.ErrorDuplicateURL, "%s %s", userErr.ErrorDuplicateURL, shortening.OriginalURL)
	}

	r.m.Store(shortening.ShortURL, shortening)

	logger.GetLogger().Debug("InMemory.Put Store Success")
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

// GetByURL метод позволяет получить запись по оригинальному URL
func (r *InMemory) GetByURL(ctx context.Context, originalURL string) (*dto.Shortening, error) {
	item := dto.Shortening{}
	r.m.Range(func(k, v interface{}) bool {
		itemCurrent := v.(dto.Shortening)
		if itemCurrent.OriginalURL == originalURL {
			item = itemCurrent
			return false
		}
		return true
	})
	return &item, nil
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

// GetStats метод для получения статистики
func (r *InMemory) GetStats(ctx context.Context) (*dto.Stats, error) {
	var (
		urls  map[string]string
		users map[string]string
	)
	r.m.Range(func(k, v interface{}) bool {
		itemCurrent := v.(dto.Shortening)

		if _, isURLExists := urls[itemCurrent.ShortURL]; !isURLExists {
			urls[itemCurrent.OriginalURL] = itemCurrent.OriginalURL
		}
		if _, isUsersExists := users[itemCurrent.ShortURL]; !isUsersExists {
			users[itemCurrent.UserID] = itemCurrent.UserID
		}
		return true
	})

	return &dto.Stats{Users: len(users), URLS: len(urls)}, nil
}
