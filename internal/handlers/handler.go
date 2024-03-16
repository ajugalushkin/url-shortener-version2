package save

import (
	"fmt"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
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
		return context.String(http.StatusBadRequest, "Wrong type request")
	}

	body, err := io.ReadAll(context.Request().Body)
	if err != nil {
		return context.String(http.StatusBadRequest, "URL parse error")
	}

	originalURL := string(body)
	if originalURL == "" {
		return context.String(http.StatusBadRequest, "URL parameter is missing")
	}

	shortenURL, err := s.servAPI.Shorten(model.ShortenInput{RawURL: originalURL})
	if err != nil {
		return context.String(http.StatusBadRequest, "URL not shortening")
	}

	shortenedURL := fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortenURL.Key)

	context.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
	context.Response().Status = http.StatusCreated
	context.Response().Write([]byte(shortenedURL))

	return context.String(http.StatusCreated, "")
}

func (s Handler) HandleRedirect(context echo.Context) error {
	if context.Request().Method != http.MethodGet {
		return context.String(http.StatusBadRequest, "Wrong type request")
	}

	key := strings.Replace(context.Request().URL.Path, "/", "", -1)

	redirect, err := s.servAPI.Redirect(key)
	if err != nil {
		return context.String(http.StatusBadRequest, "Original URL not found!")
	}

	context.Response().Header().Set("Location", redirect)
	context.Response().Status = http.StatusTemporaryRedirect
	return context.String(http.StatusTemporaryRedirect, "")
}
