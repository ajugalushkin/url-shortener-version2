package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/parse"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/validate"
	"github.com/labstack/echo/v4"
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
	if echoCtx.Request().Method != http.MethodPost {
		return validate.AddError(s.ctx, echoCtx, validate.WrongTypeRequest, http.StatusBadRequest, 0)
	}

	parseURL, err := parse.GetURL(s.ctx, echoCtx)
	if err != nil {
		return err
	}

	shortenURL, err := s.servAPI.Shorten(dto.ShortenInput{RawURL: parseURL})
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.URLNotShortening, http.StatusBadRequest, 0)
	}

	flags := config.ConfigFromContext(s.ctx)
	newBody := []byte(fmt.Sprintf("%s/%s", flags.BaseURL, shortenURL.Key))

	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)

	echoCtx.Response().Status = http.StatusCreated
	sizeBody, err := echoCtx.Response().Write(newBody)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.FailedToSend, http.StatusBadRequest, 0)
	}

	return validate.AddMessageOK(s.ctx, echoCtx, validate.URLSent, http.StatusTemporaryRedirect, sizeBody)
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
	if echoCtx.Request().Method != http.MethodPost {
		return validate.AddError(s.ctx, echoCtx, validate.WrongTypeRequest, http.StatusBadRequest, 0)
	}

	parseURL, err := parse.GetURL(s.ctx, echoCtx)
	if err != nil {
		return err
	}

	shortenURL, err := s.servAPI.Shorten(dto.ShortenInput{RawURL: parseURL})
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.URLNotShortening, http.StatusBadRequest, 0)
	}

	flags := config.ConfigFromContext(s.ctx)
	shortenResult := dto.ShortenResult{Result: fmt.Sprintf("%s/%s", flags.BaseURL, shortenURL.Key)}
	json, err := shortenResult.MarshalJSON()
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.JSONNotCreate, http.StatusBadRequest, 0)
	}

	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx.Response().Status = http.StatusCreated
	sizeBody, err := echoCtx.Response().Write(json)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.FailedToSend, http.StatusBadRequest, 0)
	}

	return validate.AddMessageOK(s.ctx, echoCtx, validate.URLSent, http.StatusTemporaryRedirect, sizeBody)
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
	if echoCtx.Request().Method != http.MethodGet {
		return validate.AddError(s.ctx, echoCtx, validate.WrongTypeRequest, http.StatusBadRequest, 0)
	}

	redirect, err := s.servAPI.Redirect(strings.Replace(echoCtx.Request().URL.Path, "/", "", -1))
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.URLNotFound, http.StatusBadRequest, 0)
	}

	echoCtx.Response().Header().Set(echo.HeaderLocation, redirect)
	echoCtx.Response().Status = http.StatusTemporaryRedirect

	return validate.AddMessageOK(s.ctx, echoCtx, validate.URLSent, http.StatusTemporaryRedirect, 0)
}
