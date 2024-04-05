package validate

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

const (
	WrongTypeRequest = "Wrong type request"
	UrlParseError    = "URL parse error"
	UrlMissing       = "URL parameter is missing"
	UrlNotShortening = "URL not shortening"
	FailedToSend     = "Failed to send URL"
	JsonParseError   = "JSON parse error"
	JsonNotCreate    = "JSON not create"
	UrlSent          = "URL sent"
	UrlNotFound      = "Original URL not found!"

	Status = "status"
	Size   = "size"
)

func CheckMethodType(ctx context.Context, context echo.Context) error {
	if context.Request().Method != http.MethodPost {
		log := logger.LoggerFromContext(ctx)
		log.Debug(WrongTypeRequest,
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, WrongTypeRequest)
	}
	return nil
}

func CheckUrlEmpty(ctx context.Context, context echo.Context) error {
	originalURL := ctx.Value("URL").(string)
	if originalURL == "" {
		log.Debug(UrlMissing,
			zap.String("status", strconv.Itoa(http.StatusBadRequest)),
			zap.String("size", strconv.Itoa(0)),
		)
		return context.String(http.StatusBadRequest, UrlMissing)
	}
	return nil
}
