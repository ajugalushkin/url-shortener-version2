package handler

import (
	"encoding/json"
	"github.com/ajugalushkin/url-shortener-version2/internal/cookies"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/validate"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

var userID int

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
		return validate.AddError(s.ctx, echoCtx, "", http.StatusUnauthorized, 0)
	}

	userID = cookies.GetUserID(s.ctx, cookieIn)
	if userID == 0 {
		return validate.AddError(s.ctx, echoCtx, "", http.StatusUnauthorized, 0)
	}

	var URLs []string
	err = json.Unmarshal(body, &URLs)
	if err != nil {
		return err
	}

	inputCh := generator(URLs)
	s.deleteURL(inputCh)

	return validate.AddMessageOK(s.ctx, echoCtx, "", http.StatusAccepted, 0)
}

func generator(input []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			inputCh <- data
		}
	}()

	return inputCh
}

func (s Handler) deleteURL(inputCh chan string) {
	log := logger.LogFromContext(s.ctx)
	for URL := range inputCh {
		err := s.servAPI.DeleteUserURL(s.ctx, URL, userID)
		if err != nil {
			log.Error(err.Error())
		}
	}
}
