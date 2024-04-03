package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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

// @Summary Shorten
// @Description Short URL
// @ID shorten
// @Accept text/plain
// @Produce text/plain
// @Success 201 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router / [post]
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

	originalURL := string(body)
	if originalURL == "" {
		logger.Log.Debug("URL parameter is missing",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "URL parameter is missing")
	}

	shortenURL, err := s.servAPI.Shorten(dto.ShortenInput{RawURL: originalURL})
	if err != nil {
		logger.Log.Debug("URL not shortening",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "URL not shortening")
	}

	newBody := []byte(fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortenURL.Key))
	context.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)

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

// @Summary ShortenJSON
// @Description Short URL in json format
// @ID shorten-json
// @Accept json
// @Produce json
// @Param input body model.Shorten true "URL for shortening"
// @Success 201 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router /api/shorten [post]
func (s Handler) HandleShorten(context echo.Context) error {
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

	shorten := dto.Shorten{}
	err = shorten.UnmarshalJSON(body)
	if err != nil {
		logger.Log.Debug("JSON parse error",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "JSON parse error")
	}

	shortenURL, err := s.servAPI.Shorten(dto.ShortenInput{RawURL: shorten.URL})
	if err != nil {
		logger.Log.Debug("URL not shortening",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "URL not shortening")
	}

	shortenResult := dto.ShortenResult{Result: fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortenURL.Key)}
	json, err := shortenResult.MarshalJSON()
	if err != nil {
		logger.Log.Debug("JSON not create",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "JSON not create")
	}

	context.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	context.Response().Status = http.StatusCreated
	_, err = context.Response().Write(json)
	if err != nil {
		logger.Log.Debug("Failed to send URL",
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, "Failed to send URL")
	}

	logger.Log.Debug("URL sent",
		zap.String("status", strconv.Itoa(http.StatusCreated)),
		zap.String("size", strconv.Itoa(len(json))),
	)

	return context.String(http.StatusTemporaryRedirect, "")
}

// @Summary Redirect
// @Description Redirect to origin URL by short URL
// @ID redirect
// @Accept text/plain
// @Produce text/html; charset=utf-8
// @Success 307 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router / [get]
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
