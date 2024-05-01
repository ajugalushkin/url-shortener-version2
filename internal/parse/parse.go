package parse

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
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

func SetResponse(ctx context.Context, echoCtx echo.Context, parseURL string, httpStatus int) error {
	contentType := echoCtx.Request().Header.Get(echo.HeaderContentType)

	echoCtx.Response().Header().Set(echo.HeaderContentType, contentType)
	echoCtx.Response().Status = httpStatus

	var newBody []byte
	if contentType != echo.MIMEApplicationJSON {
		newBody = []byte(parseURL)
	} else {
		shortenResult := dto.ShortenOutput{Result: parseURL}
		newBody, _ = shortenResult.MarshalJSON()
	}

	sizeBody, err := echoCtx.Response().Write(newBody)
	if err != nil {
		return validate.AddError(ctx, echoCtx, validate.FailedToSend, http.StatusBadRequest, 0)
	}
	return validate.AddMessageOK(ctx, echoCtx, validate.URLSent, httpStatus, sizeBody)
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
		return newBody, validate.AddError(ctx, echoCtx, validate.JSONNotCreate, http.StatusBadRequest, 0)
	}

	return newBody, nil
}

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
		return newBody, validate.AddError(ctx, echoCtx, validate.JSONNotCreate, http.StatusBadRequest, 0)
	}

	return newBody, nil
}
