package parse

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/validate"
)

// GetJSONDataFromBatch получение данных из контекста.
func GetJSONDataFromBatch(ctx context.Context, echoCtx echo.Context) (dto.ShortenListInput, error) {
	var shortList dto.ShortenListInput

	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return shortList, echoCtx.String(http.StatusBadRequest, validate.URLParseError)
	}

	err = shortList.UnmarshalJSON(body)
	if err != nil {
		return shortList, echoCtx.String(http.StatusBadRequest, validate.JSONParseError)
	}

	return shortList, nil
}

// SetJSONDataToBody внесение данных в контекст.
func SetJSONDataToBody(ctx context.Context, echoCtx echo.Context, list *dto.ShorteningList) ([]byte, error) {
	var shortenListOut dto.ShortenListOutput
	flag := config.FlagsFromContext(ctx)
	for _, item := range *list {
		shortWithHost, _ := url.JoinPath(flag.BaseURL, item.ShortURL)
		shortenListOut = append(
			shortenListOut,
			dto.ShortenListOutputLine{
				CorrelationID: item.CorrelationID,
				ShortURL:      shortWithHost,
			},
		)
	}

	newBody, err := shortenListOut.MarshalJSON()
	if err != nil {
		return newBody, echoCtx.String(http.StatusBadRequest, validate.JSONNotCreate)
	}

	return newBody, nil
}

// SetUserURLSToBody внесение данных в контекст.
func SetUserURLSToBody(ctx context.Context, echoCtx echo.Context, list *dto.ShorteningList) ([]byte, error) {
	var shortenListOut dto.UserURLList
	flag := config.FlagsFromContext(ctx)
	for _, item := range *list {
		shortWithHost, _ := url.JoinPath(flag.BaseURL, item.ShortURL)
		shortenListOut = append(
			shortenListOut,
			dto.UserURLListLine{
				ShortURL:    shortWithHost,
				OriginalURL: item.OriginalURL,
			},
		)
	}

	newBody, err := shortenListOut.MarshalJSON()
	if err != nil {
		return newBody, echoCtx.String(http.StatusBadRequest, validate.JSONNotCreate)
	}

	return newBody, nil
}
