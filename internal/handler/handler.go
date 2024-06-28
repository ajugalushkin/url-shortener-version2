package handler

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/cookies"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	userErr "github.com/ajugalushkin/url-shortener-version2/internal/errors"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/parse"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/validate"
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

// HandleSave @Summary Shorten
// @Description Short URL
// @ID shorten
// @Accept text/plain
// @Produce text/plain
// @Success 201 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router / [post]
func (s Handler) HandleSave(echoCtx echo.Context) error {
	parseURL, err := parse.GetURL(s.ctx, echoCtx)
	if err != nil {
		return echoCtx.String(http.StatusBadRequest, err.Error())
	}

	cookieValue, err := cookies.Read(echoCtx, "user")
	if err != nil {
		cookieValue = cookies.Write(s.ctx, echoCtx, "user")
	}

	shortenURL, err := s.servAPI.Shorten(s.ctx, dto.Shortening{
		OriginalURL: parseURL,
		UserID:      strconv.Itoa(cookies.GetUserID(s.ctx, cookieValue))})

	if err != nil {
		if errors.Is(err, userErr.ErrorDuplicateURL) {
			return echoCtx.String(http.StatusConflict, shortenURL.ShortURL)
		}
		return echoCtx.String(http.StatusBadRequest, err.Error())
	}

	return echoCtx.String(http.StatusCreated, shortenURL.ShortURL)
}

// HandleShorten @Summary ShortenJSON
// @Description Short URL in json format
// @ID shorten-json
// @Accept json
// @Produce json
// @Param input body dto.ShortenInput true "URL for shortening"
// @Success 201 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router /api/shorten [post]
func (s Handler) HandleShorten(echoCtx echo.Context) error {
	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return echoCtx.String(http.StatusBadRequest, validate.URLParseError)
	}

	shorten := dto.ShortenInput{}
	err = shorten.UnmarshalJSON(body)
	if err != nil || shorten.URL == "" {
		return echoCtx.String(http.StatusBadRequest, validate.JSONParseError)
	}

	shortenURL, err := s.servAPI.Shorten(s.ctx, dto.Shortening{OriginalURL: shorten.URL})
	if err != nil {
		if errors.Is(err, userErr.ErrorDuplicateURL) {
			return echoCtx.JSON(http.StatusConflict, dto.ShortenOutput{Result: shortenURL.ShortURL})
		}
		return echoCtx.String(http.StatusBadRequest, err.Error())
	}

	return echoCtx.JSON(http.StatusCreated, dto.ShortenOutput{Result: shortenURL.ShortURL})
}

// HandleShortenBatch ( @Summary ShortenBatch
// @Description Short list of URLs in json format
// @ID shorten-batch-json
// @Accept json
// @Produce json
// @Param input body dto.ShortenListInput true "URL for shortening"
// @Success 307 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router /api/shorten/batch [post]
func (s Handler) HandleShortenBatch(echoCtx echo.Context) error {
	if echoCtx.Request().Method != http.MethodPost {
		return validate.AddError(s.ctx, echoCtx, validate.WrongTypeRequest, http.StatusBadRequest, 0)
	}

	if ctType := echoCtx.Request().Header.Get(echo.HeaderContentType); ctType != echo.MIMEApplicationJSON {
		return validate.AddError(s.ctx, echoCtx, validate.WrongTypeRequest, http.StatusBadRequest, 0)
	}

	inputList, err := parse.GetJSONDataFromBatch(s.ctx, echoCtx)
	if err != nil {
		return err
	}

	shortList, err := s.servAPI.ShortenList(s.ctx, inputList)
	if err != nil {
		return err
	}

	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	echoCtx.Response().Status = http.StatusCreated

	body, err := parse.SetJSONDataToBody(s.ctx, echoCtx, shortList)
	if err != nil {
		return err
	}

	sizeBody, err := echoCtx.Response().Write(body)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.FailedToSend, http.StatusBadRequest, 0)
	}
	return validate.AddMessageOK(s.ctx, echoCtx, validate.URLSent, http.StatusTemporaryRedirect, sizeBody)
}

// HandleRedirect @Summary Redirect
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

	redirect, err := s.servAPI.Redirect(s.ctx, strings.Replace(echoCtx.Request().URL.Path, "/", "", -1))
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.URLNotFound, http.StatusBadRequest, 0)
	}

	if redirect.IsDeleted {
		return validate.AddError(s.ctx, echoCtx, "URL was delete!", http.StatusGone, 0)
	}

	if redirect.OriginalURL != "" {
		return validate.Redirect(s.ctx, echoCtx, redirect.OriginalURL)
	}

	log := logger.LogFromContext(s.ctx)
	log.Error(validate.URLNotFound)
	return validate.AddError(s.ctx, echoCtx, validate.URLNotFound, http.StatusBadRequest, 0)
}

// HandlePing ( @Summary Ping
// @Description Ping Database for check connection
// @ID ping
// @Accept text/plain
// @Produce text/plain; charset=utf-8
// @Success 200 {integer} integer 1
// @Failure 500 {integer} integer 1
// @Router /ping [get]
func (s Handler) HandlePing(echoCtx echo.Context) error {
	flags := config.FlagsFromContext(s.ctx)
	db, err := sql.Open("pgx", flags.DataBaseDsn)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, "", http.StatusInternalServerError, 0)
	}
	defer db.Close()

	return validate.AddMessageOK(s.ctx, echoCtx, "", http.StatusOK, 0)
}

// HandleUserUrls ( @Summary UserURLS
// @Description Retrive all short URLS for user
// @ID user-urls-json
// @Accept json
// @Produce json
// @Success 307 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router /api/user/urls [get]
func (s Handler) HandleUserUrls(echoCtx echo.Context) error {
	if echoCtx.Request().Method != http.MethodGet {
		return validate.AddError(s.ctx, echoCtx, validate.WrongTypeRequest, http.StatusBadRequest, 0)
	}

	cookieIn, err := cookies.Read(echoCtx, "user")
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, "", http.StatusUnauthorized, 0)
	}

	userID := cookies.GetUserID(s.ctx, cookieIn)
	if userID == 0 {
		return validate.AddError(s.ctx, echoCtx, "", http.StatusUnauthorized, 0)
	}

	shortList, err := s.servAPI.GetUserURLS(s.ctx, userID)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.URLNotFound, http.StatusBadRequest, 0)
	}

	body, err := parse.SetUserURLSToBody(s.ctx, echoCtx, shortList)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, "", http.StatusNoContent, 0)
	}

	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	sizeBody, err := echoCtx.Response().Write(body)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.FailedToSend, http.StatusBadRequest, 0)
	}
	return validate.AddMessageOK(s.ctx, echoCtx, validate.URLSent, http.StatusTemporaryRedirect, sizeBody)
}

// HandleUserUrlsDelete ( @Summary UserURLSDelete
// @Description Delete all short URLS for user
// @ID user-urls-json
// @Accept json
// @Produce json
// @Success 202 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router /api/user/urls [delete]
func (s Handler) HandleUserUrlsDelete(echoCtx echo.Context) error {
	if echoCtx.Request().Method != http.MethodDelete {
		return validate.AddError(s.ctx, echoCtx, validate.WrongTypeRequest, http.StatusBadRequest, 0)
	}

	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.URLParseError, http.StatusBadRequest, 0)
	}

	cookieIn, err := cookies.Read(echoCtx, "user")
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, "Not found cookies fo user", http.StatusUnauthorized, 0)
	}

	userID := cookies.GetUserID(s.ctx, cookieIn)
	if userID == 0 {
		return validate.AddError(s.ctx, echoCtx, "Not found UserID for user", http.StatusUnauthorized, 0)
	}

	var URLs dto.URLs
	err = URLs.UnmarshalJSON(body)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, "Error parse input json", http.StatusUnauthorized, 0)
	}

	s.servAPI.DeleteUserURL(s.ctx, URLs, userID)

	return validate.AddMessageOK(s.ctx, echoCtx, "URLS Delete OK", http.StatusAccepted, 0)
}
