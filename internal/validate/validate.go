package validate

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
)

// блок констант
const (
	WrongTypeRequest = "Wrong type request"
	URLParseError    = "URL parse error"
	URLMissing       = "URL parameter is missing"
	FailedToSend     = "Failed to send URL"
	JSONParseError   = "JSON parse error"
	JSONNotCreate    = "JSON not create"
	URLSent          = "URL sent"
	URLNotFound      = "Original URL not found!"

	Status = "status"
	Size   = "size"
)

// AddError функция добавления ошибки
func AddError(ctx context.Context, echoCtx echo.Context, message string, httpStatus int, size int) error {
	log := logger.LogFromContext(ctx)
	log.Debug(message,
		zap.String(Status, strconv.Itoa(httpStatus)),
		zap.String(Size, strconv.Itoa(size)),
	)

	return echoCtx.String(httpStatus, message)
}

// AddMessageOK функция добавления сообщения
func AddMessageOK(ctx context.Context, echoCtx echo.Context, message string, httpStatus int, size int) error {
	log := logger.LogFromContext(ctx)
	log.Debug(message,
		zap.String(Status, strconv.Itoa(httpStatus)),
		zap.String(Size, strconv.Itoa(size)),
	)

	return echoCtx.String(httpStatus, "")
}

// Redirect функция редиректа
func Redirect(ctx context.Context, echoCtx echo.Context, redirect string) error {
	log := logger.LogFromContext(ctx)
	log.Debug(URLSent,
		zap.String(Status, strconv.Itoa(http.StatusTemporaryRedirect)),
		zap.String(Size, strconv.Itoa(len(redirect))),
	)

	return echoCtx.Redirect(http.StatusTemporaryRedirect, redirect)
}
