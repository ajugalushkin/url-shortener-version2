package parse

import (
	"context"
	"io"
	"net/http"

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
		shorten := dto.Shorten{}
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
