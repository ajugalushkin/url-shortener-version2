package handler

import (
	"fmt"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	servAPI *service.Service
	cfg     *config.Config
}

func NewHandler(servAPI *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		servAPI: servAPI,
		cfg:     cfg}
}

func (s Handler) HandleSave(context echo.Context) error {
	if context.Request().Method != http.MethodPost {
		logger.Log.Debug("Wrong type request",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "Wrong type request")
	}

	body, err := io.ReadAll(context.Request().Body)
	if err != nil {
		logger.Log.Debug("URL parse error",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "URL parse error")
	}

	contentType := context.Request().Header.Get("Content-Type")

	var originalURL string
	if contentType == echo.MIMEApplicationJSON {
		shorten := model.Shorten{}
		err = shorten.UnmarshalJSON(body)
		if err != nil {
			logger.Log.Debug("JSON parse error",
				zap.String("status", strconv.Itoa(http.StatusBadRequest)),
				zap.String("size", strconv.Itoa(0)),
			)
			return context.String(http.StatusBadRequest, "JSON parse error")
		}
		originalURL = shorten.URL
	} else {
		originalURL = string(body)
	}

	if originalURL == "" {
		logger.Log.Debug("URL parameter is missing",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "URL parameter is missing")
	}

	shortenURL, err := s.servAPI.Shorten(model.ShortenInput{RawURL: originalURL})
	if err != nil {
		logger.Log.Debug("URL not shortening",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "URL not shortening")
	}

	var newBody []byte
	if contentType == echo.MIMEApplicationJSON {
		shortenResult := model.ShortenResult{Result: fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortenURL.Key)}
		newBody, err = shortenResult.MarshalJSON()
		if err != nil {
			logger.Log.Debug("JSON not create",
				zap.String("status", strconv.Itoa(http.StatusBadRequest)),
				zap.String("size", strconv.Itoa(0)),
			)
			return context.String(http.StatusBadRequest, "JSON not create")
		}

		context.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		newBody = []byte(fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortenURL.Key))
		context.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
	}

	context.Response().Status = http.StatusCreated
	_, err = context.Response().Write(newBody)
	if err != nil {
		logger.Log.Debug("Failed to send URL",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "Failed to send URL")
	}

	logger.Log.Debug("URL sent",
		zap.String("status", strconv.Itoa(http.StatusCreated)),
		zap.String("size", strconv.Itoa(len(newBody))),
	)

	return context.String(http.StatusCreated, "")
}

func (s Handler) HandleRedirect(context echo.Context) error {
	if context.Request().Method != http.MethodGet {
		logger.Log.Debug("Wrong type request",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "Wrong type request")
	}

	key := strings.Replace(context.Request().URL.Path, "/", "", -1)

	redirect, err := s.servAPI.Redirect(key)
	if err != nil {
		logger.Log.Debug("Original URL not found!",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "Original URL not found!")
	}

	context.Response().Header().Set("Location", redirect)
	context.Response().Status = http.StatusTemporaryRedirect

	logger.Log.Debug("Original URL not found!",
		zap.String("status", strconv.Itoa(http.StatusTemporaryRedirect)),
		zap.String("size", strconv.Itoa(0)),
	)

	return context.String(http.StatusTemporaryRedirect, "")
}
