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

// Handler структура Handler
type Handler struct {
	ctx     context.Context
	cache   map[string]*dto.User
	servAPI *service.Service
}

// NewHandler конструктор
func NewHandler(ctx context.Context, servAPI *service.Service) *Handler {
	return &Handler{
		ctx:     ctx,
		cache:   make(map[string]*dto.User),
		servAPI: servAPI}
}

const cookieName string = "User"

// HandleSave @Summary Shorten
// @Description Short URL
// @ID shorten
// @Accept text/plain
// @Produce text/plain
// @Success 201 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router / [post]
func (s Handler) HandleSave(echoCtx echo.Context) error {
	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil || len(body) == 0 {
		return echoCtx.String(http.StatusBadRequest, "Body is empty or invalid")
	}

	cookieValue, err := cookies.Read(echoCtx, cookieName)
	if err != nil {
		cookieValue = cookies.Write(s.ctx, echoCtx, cookieName)
		s.cache[cookieValue] = cookies.GetUser(s.ctx, cookieValue)
	}

	shortenURL, err := s.servAPI.Shorten(s.ctx, dto.Shortening{
		OriginalURL: string(body),
		UserID:      strconv.Itoa(cookies.GetUser(s.ctx, cookieValue).ID)})

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
	if ctType := echoCtx.Request().Header.Get(echo.HeaderContentType); ctType != echo.MIMEApplicationJSON {
		return echoCtx.String(http.StatusBadRequest, validate.WrongTypeRequest)
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

	_, err = echoCtx.Response().Write(body)
	if err != nil {
		return echoCtx.String(http.StatusBadRequest, validate.FailedToSend)
	}
	return echoCtx.String(http.StatusTemporaryRedirect, "")
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
	redirect, err := s.servAPI.Redirect(s.ctx, strings.Replace(echoCtx.Request().URL.Path, "/", "", -1))
	if err != nil {
		return echoCtx.String(http.StatusBadRequest, validate.URLNotFound)
	}

	if redirect.IsDeleted {
		return echoCtx.String(http.StatusGone, validate.URLWasDelete)
	}

	if redirect.OriginalURL != "" {
		return echoCtx.Redirect(http.StatusTemporaryRedirect, redirect.OriginalURL)
	}

	logger.GetLogger().Error(validate.URLNotFound)
	return echoCtx.String(http.StatusBadRequest, validate.URLNotFound)
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
	flags := config.GetConfig()
	db, err := sql.Open("pgx", flags.DataBaseDsn)
	if err != nil {
		return echoCtx.String(http.StatusInternalServerError, validate.PingError)
	}
	defer db.Close()

	return echoCtx.String(http.StatusOK, validate.PingOk)
}

// CustomContext структура расширяет echo.Context
type CustomContext struct {
	user *dto.User
	echo.Context
}

// Authorized middleware для авторизация cookie
func (s Handler) Authorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		cookie, err := echoCtx.Cookie(cookieName)
		if err != nil {
			return echoCtx.String(http.StatusUnauthorized, err.Error())
		}

		if _, ok := s.cache[cookie.Value]; !ok {
			return echoCtx.String(http.StatusUnauthorized, "")
		}

		user := cookies.GetUser(s.ctx, cookie.Value)

		newContext := &CustomContext{user: user, Context: echoCtx}

		return next(newContext)
	}
}

// HandleUserUrls ( @Summary UserURLS
// @Description Retrive all short URLS for user
// @ID user-urls-json
// @Accept json
// @Produce json
// @Success 307 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router /api/user/urls [get]
func (s Handler) HandleUserUrls(c echo.Context) error {
	echoCtx := c.(*CustomContext)

	shortList, err := s.servAPI.GetUserURLS(s.ctx, echoCtx.user.ID)
	if err != nil {
		return echoCtx.String(http.StatusBadRequest, validate.URLNotFound)
	}

	body, err := parse.SetUserURLSToBody(s.ctx, echoCtx, shortList)
	if err != nil {
		return echoCtx.String(http.StatusNoContent, "")
	}

	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	_, err = echoCtx.Response().Write(body)
	if err != nil {
		return echoCtx.String(http.StatusBadRequest, validate.FailedToSend)
	}
	return echoCtx.String(http.StatusTemporaryRedirect, "")
}

// HandleUserUrlsDelete ( @Summary UserURLSDelete
// @Description Delete all short URLS for user
// @ID user-urls-json
// @Accept json
// @Produce json
// @Success 202 {integer} integer 1
// @Failure 400 {integer} integer 1
// @Router /api/user/urls [delete]
func (s Handler) HandleUserUrlsDelete(c echo.Context) error {
	echoCtx := c.(*CustomContext)

	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return echoCtx.String(http.StatusBadRequest, validate.URLParseError)
	}

	var URLs dto.URLs
	err = URLs.UnmarshalJSON(body)
	if err != nil {
		return echoCtx.String(http.StatusInternalServerError, "Error parse input json")
	}

	s.servAPI.DeleteUserURL(s.ctx, URLs, echoCtx.user.ID)

	return echoCtx.String(http.StatusAccepted, "URLS Delete OK")
}
