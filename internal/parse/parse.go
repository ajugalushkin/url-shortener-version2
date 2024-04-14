package parse

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/validate"
	"github.com/labstack/echo/v4"
)

func GetURL(ctx context.Context, echoCtx echo.Context) (string, error) {
	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return "", validate.AddError(ctx, echoCtx, validate.URLParseError, http.StatusBadRequest, 0)
	}

	var parseURL string
	contentType := echoCtx.Request().Header.Get(echo.HeaderContentType)
	if contentType != echo.MIMEApplicationJSON {
		parseURL = string(body)
	} else {
		shorten := dto.ShortenInput{}
		err = shorten.UnmarshalJSON(body)
		if err != nil {
			return "", validate.AddError(ctx, echoCtx, validate.JSONParseError, http.StatusBadRequest, 0)
		}
		parseURL = shorten.URL
	}

	if parseURL == "" {
		return "", validate.AddError(ctx, echoCtx, validate.URLMissing, http.StatusBadRequest, 0)
	}

	return parseURL, nil
}
func SetBody(ctx context.Context, echoCtx echo.Context, servAPI *service.Service, parseURL string) ([]byte, error) {
	var newBody []byte
	shortenURL, err := servAPI.Shorten(dto.Shortening{OriginalURL: parseURL})
	if err != nil {
		return newBody, validate.AddError(ctx, echoCtx, validate.URLNotShortening, http.StatusBadRequest, 0)
	}

	flags := config.FlagsFromContext(ctx)

	contentType := echoCtx.Request().Header.Get(echo.HeaderContentType)
	if contentType != echo.MIMEApplicationJSON {
		newBody = []byte(fmt.Sprintf("%s/%s", flags.BaseURL, shortenURL.ShortURL))
	} else {
		shortenResult := dto.ShortenOutput{Result: fmt.Sprintf("%s/%s", flags.BaseURL, shortenURL.ShortURL)}
		newBody, err = shortenResult.MarshalJSON()
		if err != nil {
			return newBody, validate.AddError(ctx, echoCtx, validate.JSONNotCreate, http.StatusBadRequest, 0)
		}
	}
	return newBody, nil
}

func GetJSONDataFromBatch(ctx context.Context, echoCtx echo.Context) (dto.ShortenListInput, error) {
	var shortList dto.ShortenListInput

	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return shortList, validate.AddError(ctx, echoCtx, validate.URLParseError, http.StatusBadRequest, 0)
	}

	err = shortList.UnmarshalJSON(body)
	if err != nil {
		return shortList, validate.AddError(ctx, echoCtx, validate.JSONParseError, http.StatusBadRequest, 0)
	}

	return shortList, nil
}

func SetJSONDataToBody(ctx context.Context, echoCtx echo.Context, list *dto.ShorteningList) ([]byte, error) {
	var shortenListOut dto.ShortenListOutput
	for _, item := range *list {
		shortenListOut = append(
			shortenListOut,
			dto.ShortenListOutputLine{
				CorrelationID: item.CorrelationID,
				ShortURL:      item.ShortURL,
			},
		)
	}

	newBody, err := shortenListOut.MarshalJSON()
	if err != nil {
		return newBody, validate.AddError(ctx, echoCtx, validate.JSONNotCreate, http.StatusBadRequest, 0)
	}

	return newBody, nil
}
