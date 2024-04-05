package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/validate"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	ctx     context.Context
	servAPI *service.Service
}

func NewHandler(ctx context.Context, servAPI *service.Service) *Handler {
	return &Handler{
		ctx:     ctx,
		servAPI: servAPI}
}

// @Summary Shorten
// @Description Short URL
// @ID shorten
// @Accept text/plain
// @Produce text/plain
// @Success 201 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router / [post]
func (s Handler) HandleSave(echoCtx echo.Context) error {
	log := logger.LoggerFromContext(s.ctx)

	err := validate.CheckMethodType(s.ctx, echoCtx)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		log.Debug(validate.UrlParseError,
			zap.String(validate.Status, strconv.Itoa(http.StatusBadRequest)),
			zap.String(validate.Size, strconv.Itoa(0)),
		)
		return echoCtx.String(http.StatusBadRequest, validate.UrlParseError)
	}

	originalURL := string(body)

	s.ctx = context.WithValue(s.ctx, "URL", originalURL)
	err = validate.CheckUrlEmpty(s.ctx, echoCtx)
	if err != nil {
		return err
	}

	shortenURL, err := s.servAPI.Shorten(dto.ShortenInput{RawURL: originalURL})
	if err != nil {
		log.Debug(validate.UrlNotShortening,
			zap.String(validate.Status, strconv.Itoa(http.StatusBadRequest)),
			zap.String(validate.Size, strconv.Itoa(0)),
		)
		return echoCtx.String(http.StatusBadRequest, validate.UrlNotShortening)
	}

	flags := config.ConfigFromContext(s.ctx)
	newBody := []byte(fmt.Sprintf("%s/%s", flags.BaseURL, shortenURL.Key))
	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)

	echoCtx.Response().Status = http.StatusCreated
	_, err = echoCtx.Response().Write(newBody)
	if err != nil {
		log.Debug(validate.FailedToSend,
			zap.String(validate.Status, strconv.Itoa(http.StatusBadRequest)),
			zap.String(validate.Size, strconv.Itoa(0)),
		)
		return echoCtx.String(http.StatusBadRequest, validate.FailedToSend)
	}

	log.Debug(validate.UrlSent,
		zap.String(validate.Status, strconv.Itoa(http.StatusCreated)),
		zap.String(validate.Size, strconv.Itoa(len(newBody))),
	)

	return echoCtx.String(http.StatusCreated, "")
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
func (s Handler) HandleShorten(echoCtx echo.Context) error {
	log := logger.LoggerFromContext(s.ctx)

	err := validate.CheckMethodType(s.ctx, echoCtx)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		log.Debug(validate.UrlParseError,
			zap.String(validate.Status, strconv.Itoa(http.StatusBadRequest)),
			zap.String(validate.Size, strconv.Itoa(0)),
		)
		return echoCtx.String(http.StatusBadRequest, validate.UrlParseError)
	}

	shorten := dto.Shorten{}
	err = shorten.UnmarshalJSON(body)
	if err != nil {
		log.Debug(validate.JsonParseError,
			zap.String(validate.Status, strconv.Itoa(http.StatusBadRequest)),
			zap.String(validate.Size, strconv.Itoa(0)),
		)
		return echoCtx.String(http.StatusBadRequest, validate.JsonParseError)
	}

	s.ctx = context.WithValue(s.ctx, "URL", shorten.URL)
	err = validate.CheckUrlEmpty(s.ctx, echoCtx)
	if err != nil {
		return err
	}

	shortenURL, err := s.servAPI.Shorten(dto.ShortenInput{RawURL: shorten.URL})
	if err != nil {
		log.Debug(validate.UrlNotShortening,
			zap.String(validate.Status, strconv.Itoa(http.StatusBadRequest)),
			zap.String(validate.Size, strconv.Itoa(0)),
		)
		return echoCtx.String(http.StatusBadRequest, validate.UrlNotShortening)
	}

	flags := config.ConfigFromContext(s.ctx)
	shortenResult := dto.ShortenResult{Result: fmt.Sprintf("%s/%s", flags.BaseURL, shortenURL.Key)}
	json, err := shortenResult.MarshalJSON()
	if err != nil {
		log.Debug(validate.JsonNotCreate,
			zap.String(validate.Status, strconv.Itoa(http.StatusBadRequest)),
			zap.String(validate.Size, strconv.Itoa(0)),
		)
		return echoCtx.String(http.StatusBadRequest, validate.JsonNotCreate)
	}

	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx.Response().Status = http.StatusCreated
	_, err = echoCtx.Response().Write(json)
	if err != nil {
		log.Debug(validate.FailedToSend,
			zap.String(validate.Status, strconv.Itoa(http.StatusBadRequest)),
			zap.String(validate.Size, strconv.Itoa(0)),
		)
		return echoCtx.String(http.StatusBadRequest, validate.FailedToSend)
	}

	log.Debug(validate.UrlSent,
		zap.String(validate.Status, strconv.Itoa(http.StatusCreated)),
		zap.String(validate.Size, strconv.Itoa(len(json))),
	)

	return echoCtx.String(http.StatusTemporaryRedirect, "")
}

// @Summary Redirect
// @Description Redirect to origin URL by short URL
// @ID redirect
// @Accept text/plain
// @Produce text/html; charset=utf-8
// @Success 307 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router / [get]
func (s Handler) HandleRedirect(echoCtx echo.Context) error {
	log := logger.LoggerFromContext(s.ctx)

	err := validate.CheckMethodType(s.ctx, echoCtx)
	if err != nil {
		return err
	}

	key := strings.Replace(echoCtx.Request().URL.Path, "/", "", -1)

	redirect, err := s.servAPI.Redirect(key)
	if err != nil {
		log.Debug(validate.UrlNotFound,
			zap.String(validate.Status, strconv.Itoa(http.StatusBadRequest)),
			zap.String(validate.Size, strconv.Itoa(0)),
		)
		return echoCtx.String(http.StatusBadRequest, validate.UrlNotFound)
	}

	echoCtx.Response().Header().Set("Location", redirect)
	echoCtx.Response().Status = http.StatusTemporaryRedirect

	log.Debug(validate.UrlSent,
		zap.String(validate.Status, strconv.Itoa(http.StatusTemporaryRedirect)),
		zap.String(validate.Size, strconv.Itoa(0)),
	)

	return echoCtx.String(http.StatusTemporaryRedirect, "")
}
