package handler

import (
	"encoding/json"
	"github.com/ajugalushkin/url-shortener-version2/internal/cookies"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/validate"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"sync"
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

	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := s.generator(doneCh, URLs)
	channels := s.fanOut(doneCh, inputCh)
	s.fanIn(doneCh, channels...)

	return validate.AddMessageOK(s.ctx, echoCtx, "", http.StatusAccepted, 0)
}

func (s Handler) generator(doneCh chan struct{}, input []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			select {
			case <-doneCh:
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}

func (s Handler) fanOut(doneCh chan struct{}, inputCh chan string) []chan string {
	numWorkers := 10
	channels := make([]chan string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		addResultCh := s.deleteURL(doneCh, inputCh)
		channels[i] = addResultCh
	}
	return channels
}

func (s Handler) fanIn(doneCh chan struct{}, resultChs ...chan string) {
	finalCh := make(chan string)

	var wg sync.WaitGroup

	for _, ch := range resultChs {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-doneCh:
					return
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	log := logger.LogFromContext(s.ctx)
	for res := range finalCh {
		log.Info(res)
	}
}

func (s Handler) deleteURL(doneCh chan struct{}, inputCh chan string) chan string {
	addRes := make(chan string)

	go func() {
		defer close(addRes)

		for data := range inputCh {
			err := s.servAPI.DeleteUserURL(s.ctx, data, userID)
			var result = "DeleteUserURL OK"
			if err != nil {
				result = err.Error()
			}

			select {
			case <-doneCh:
				return
			case addRes <- result:
			}
		}
	}()
	return addRes
}
