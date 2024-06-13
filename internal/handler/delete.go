package handler

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/cookies"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/validate"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func (s Handler) HandleUserUrlsDelete(echoCtx echo.Context) error {
	if echoCtx.Request().Method != http.MethodDelete {
		return validate.AddError(s.ctx, echoCtx, validate.WrongTypeRequest, http.StatusBadRequest, 0)
	}

	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, validate.URLParseError, http.StatusBadRequest, 0)
	}

	cookieIn, err := cookies.Read(echoCtx, "user")
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, "Not found cookies fo user", http.StatusUnauthorized, 0)
	}

	userID := cookies.GetUserID(s.ctx, cookieIn)
	if userID == 0 {
		return validate.AddError(s.ctx, echoCtx, "Not found UserID for user", http.StatusUnauthorized, 0)
	}

	var URLs dto.URLs
	err = URLs.UnmarshalJSON(body)
	if err != nil {
		return validate.AddError(s.ctx, echoCtx, "Error parse input json", http.StatusUnauthorized, 0)
	}

	s.servAPI.DeleteUserURL(s.ctx, URLs, userID)

	return validate.AddMessageOK(s.ctx, echoCtx, "URLS Delete OK", http.StatusAccepted, 0)
}
