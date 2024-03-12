package redirect

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func New(serviceAPI *service.Service) echo.HandlerFunc {
	return func(context echo.Context) error {
		if context.Request().Method != http.MethodGet {
			return context.String(http.StatusBadRequest, "Wrong type request")
		}

		key := strings.Replace(context.Request().URL.Path, "/", "", -1)

		redirect, err := serviceAPI.Redirect(key)
		if err != nil {
			return context.String(http.StatusBadRequest, "Original URL not found!")
		}

		context.Response().Header().Set("Location", redirect)
		context.Response().Status = http.StatusTemporaryRedirect
		return context.String(http.StatusTemporaryRedirect, "")
	}
}
