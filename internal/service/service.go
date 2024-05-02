package service

import (
	"context"
	"net/url"
	"strconv"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/shorten"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PutGetter interface {
	Put(ctx context.Context, shortening dto.Shortening) (*dto.Shortening, error)
	Get(ctx context.Context, shortURL string) (*dto.Shortening, error)
	GetListByUser(ctx context.Context, userID string) (*dto.ShorteningList, error)
	PutList(ctx context.Context, list dto.ShorteningList) error
	DeleteUserURL(ctx context.Context, shortURL string, userID int) error
}

type Service struct {
	storage PutGetter
}

func NewService(storage PutGetter) *Service {
	return &Service{storage: storage}
}

func (s *Service) Shorten(ctx context.Context, input dto.Shortening) (*dto.Shortening, error) {
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
		CorrelationID: input.CorrelationID,
		UserID:        input.UserID,
		IsDeleted:     input.IsDeleted,
	}

	shortening, err := s.storage.Put(ctx, newShortening)
	shortening.ShortURL, _ = url.JoinPath(config.FlagsFromContext(ctx).BaseURL, shortening.ShortURL)

	if err != nil {
		return shortening, err
	}

	return shortening, nil
}

func (s *Service) ShortenList(ctx context.Context, input dto.ShortenListInput) (*dto.ShorteningList, error) {
	var shorteningList dto.ShorteningList
	for _, item := range input {
		newShortening := dto.Shortening{
			ShortURL:      shorten.Shorten(uuid.New().ID()),
			OriginalURL:   item.OriginalURL,
			CorrelationID: item.CorrelationID,
		}

		shorteningList = append(shorteningList, newShortening)
	}

	err := s.storage.PutList(ctx, shorteningList)
	if err != nil {
		return nil, err
	}

	return &shorteningList, nil
}

func (s *Service) Redirect(ctx context.Context, identifier string) (*dto.Shortening, error) {
	log := logger.LogFromContext(ctx)

	shortening, err := s.storage.Get(ctx, identifier)
	if err != nil {
		log.Info("service.Redirect ERROR", zap.Error(err))
		return shortening, err
	}

	log.Info("service.Redirect OK", zap.String("URL", shortening.OriginalURL))

	return shortening, nil
}

func (s *Service) GetUserURLS(ctx context.Context, userID int) (*dto.ShorteningList, error) {
	log := logger.LogFromContext(ctx)

	shortening, err := s.storage.GetListByUser(ctx, strconv.Itoa(userID))
	if err != nil {
		log.Info("service.GetUserURLS ERROR", zap.Error(err))
		return &dto.ShorteningList{}, err
	}
	return shortening, nil
}

func (s *Service) DeleteUserURL(ctx context.Context, shortURL string, userID int) error {
	log := logger.LogFromContext(ctx)

	err := s.storage.DeleteUserURL(ctx, shortURL, userID)
	if err != nil {
		log.Info("service.DeleteUserURLS ERROR", zap.Error(err))
		return err
	}
	return nil
}
