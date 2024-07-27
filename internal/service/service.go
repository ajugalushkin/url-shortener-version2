package service

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"go.uber.org/zap"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	userErr "github.com/ajugalushkin/url-shortener-version2/internal/errors"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/shorten"
)

// PutGetter интерфейс для хранилища
type PutGetter interface {
	Put(ctx context.Context, shortening dto.Shortening) (*dto.Shortening, error)
	Get(ctx context.Context, shortURL string) (*dto.Shortening, error)
	GetListByUser(ctx context.Context, userID string) (*dto.ShorteningList, error)
	PutList(ctx context.Context, list dto.ShorteningList) error
	DeleteUserURL(ctx context.Context, shortURL []string, userID int)
}

// Service структура сервиса
type Service struct {
	storage PutGetter
}

// NewService конструктор
func NewService(storage PutGetter) *Service {
	return &Service{storage: storage}
}

// Shorten метод для получения сокращенного URL
func (s *Service) Shorten(ctx context.Context, input dto.Shortening) (*dto.Shortening, error) {
	if input.OriginalURL == "" {
		return nil, errors.New("URL is empty")
	}

	var (
		identifier = input.ShortURL
	)

	logger.GetLogger().Debug("Service.Shorten",
		zap.String("Origin URL", input.OriginalURL))

	if identifier == "" {
		identifier = shorten.Shorten(input.OriginalURL)
	}

	logger.GetLogger().Debug("Service.Shorten",
		zap.String("Short URL", identifier))

	newShortening := dto.Shortening{
		ShortURL:      identifier,
		OriginalURL:   input.OriginalURL,
		CorrelationID: input.CorrelationID,
		UserID:        input.UserID,
		IsDeleted:     input.IsDeleted,
	}

	shortening, err := s.storage.Put(ctx, newShortening)

	if errors.Is(err, userErr.ErrorDuplicateURL) || shortening != nil {
		shortening.ShortURL, _ = url.JoinPath(config.GetConfig().BaseURL, shortening.ShortURL)
	}

	if err != nil {
		logger.GetLogger().Debug("Service.Shorten Put Error",
			zap.Error(err))
		return shortening, err
	}

	logger.GetLogger().Debug("Service.Shorten Ok",
		zap.String("Shorten URL", shortening.ShortURL))
	return shortening, nil
}

// ShortenList метод для получения сокращения списка URL
func (s *Service) ShortenList(ctx context.Context, input dto.ShortenListInput) (*dto.ShorteningList, error) {
	var shorteningList dto.ShorteningList
	for _, item := range input {
		newShortening := dto.Shortening{
			ShortURL:      shorten.Shorten(item.OriginalURL),
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

// Redirect метод для перенаправления
func (s *Service) Redirect(ctx context.Context, identifier string) (*dto.Shortening, error) {
	log := logger.GetLogger()

	shortening, err := s.storage.Get(ctx, identifier)
	if err != nil {
		log.Debug("service.Redirect ERROR", zap.Error(err))
		return shortening, err
	}
	return shortening, nil
}

// GetUserURLS метод для получения списка URL для конкретного пользователя.
func (s *Service) GetUserURLS(ctx context.Context, userID int) (*dto.UserURLList, error) {
	log := logger.GetLogger()

	shorteningList, err := s.storage.GetListByUser(ctx, strconv.Itoa(userID))
	if err != nil {
		log.Info("service.GetUserURLS ERROR", zap.Error(err))
		return &dto.UserURLList{}, err
	}

	var shortenListOut dto.UserURLList
	for _, item := range *shorteningList {
		shortWithHost, _ := url.JoinPath(config.GetConfig().BaseURL, item.ShortURL)
		shortenListOut = append(
			shortenListOut,
			dto.UserURLListLine{
				ShortURL:    shortWithHost,
				OriginalURL: item.OriginalURL,
			},
		)
	}

	//newBody, err := shortenListOut.MarshalJSON()
	//if err != nil {
	//	return newBody, echoCtx.String(http.StatusBadRequest, validate.JSONNotCreate)
	//}

	return &shortenListOut, nil
}

// DeleteUserURL метод для удаления списка URL для конкретного пользователя.
func (s *Service) DeleteUserURL(ctx context.Context, shortURL []string, userID int) {
	s.storage.DeleteUserURL(ctx, shortURL, userID)
}
