package validate

import (
	"context"
	"strconv"

	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/labstack/echo/v4"
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

func AddError(ctx context.Context, echoCtx echo.Context, message string, httpStatus int, size int) error {
	logger.LoggerFromContext(ctx).Debug(message,
		zap.String(Status, strconv.Itoa(httpStatus)),
		zap.String(Size, strconv.Itoa(size)),
	)

	return echoCtx.String(httpStatus, message)
}

func AddMessageOK(ctx context.Context, echoCtx echo.Context, message string, httpStatus int, size int) error {
	logger.LoggerFromContext(ctx).Debug(message,
		zap.String(Status, strconv.Itoa(httpStatus)),
		zap.String(Size, strconv.Itoa(size)),
	)

	return echoCtx.String(httpStatus, "")
}
