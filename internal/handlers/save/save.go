package save

import (
	"fmt"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func New(serviceAPI *service.Service) echo.HandlerFunc {
	return func(context echo.Context) error {
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

		shortenURL, err := serviceAPI.Shorten(model.ShortenInput{RawURL: originalURL})
		if err != nil {
			return context.String(http.StatusBadRequest, "URL not shortening")
		}

		shortenedURL := fmt.Sprintf("%s/%s", config.BaseURL, shortenURL.Key)

		context.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
		context.Response().Status = http.StatusCreated
		context.Response().Write([]byte(shortenedURL))

		return context.String(http.StatusCreated, "")
	}
}
