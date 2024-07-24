package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-resty/resty/v2"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
)

func BenchmarkHandler_HandleSave(b *testing.B) {
	b.Run("Endpoint: POST /", func(b *testing.B) {
		e := echo.New()
		h := NewHandler(ctx, service.NewService(storage.GetStorage(ctx)))

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(gofakeit.URL()))
			req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
			rec := httptest.NewRecorder()
			newContext := e.NewContext(req, rec)

			err := h.HandleSave(newContext)
			if err != nil {
				return
			}
		}
	})
}

func BenchmarkHandler_HandleShorten(b *testing.B) {
	b.Run("Endpoint: POST /api/shorten", func(b *testing.B) {
		e := echo.New()
		h := NewHandler(ctx, service.NewService(storage.GetStorage(ctx)))

		for i := 0; i < b.N; i++ {
			shorten := dto.ShortenInput{URL: gofakeit.URL()}
			json, _ := shorten.MarshalJSON()

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/shorten",
				strings.NewReader(string(json)),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoContext := e.NewContext(req, rec)

			err := h.HandleShorten(echoContext)
			if err != nil {
				return
			}
		}
	})
}

func BenchmarkHandler_HandleShortenBatch(b *testing.B) {
	b.Run("Endpoint: POST /api/shorten/batch", func(b *testing.B) {
		e := echo.New()
		h := NewHandler(ctx, service.NewService(storage.GetStorage(ctx)))

		for i := 0; i < b.N; i++ {
			shortList := dto.ShortenListInput{{
				CorrelationID: gofakeit.UUID(),
				OriginalURL:   gofakeit.URL(),
			}}
			json, _ := shortList.MarshalJSON()

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/shorten/batch",
				strings.NewReader(string(json)),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoContext := e.NewContext(req, rec)

			err := h.HandleShorten(echoContext)
			if err != nil {
				return
			}
		}
	})
}

func BenchmarkHandler_HandleRedirect(b *testing.B) {
	b.Run("Endpoint: GET /", func(b *testing.B) {
		e := echo.New()
		h := NewHandler(ctx, service.NewService(storage.GetStorage(ctx)))

		for i := 0; i < b.N; i++ {
			// Post
			client := resty.New()

			shortenInput := dto.ShortenInput{URL: gofakeit.URL()}
			json, _ := shortenInput.MarshalJSON()

			shortenOut := dto.ShortenOutput{}

			_, err := client.R().
				SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
				SetBody(json).
				SetResult(&shortenOut).
				Post(config.GetConfig().BaseURL + "/api/shorten")
			if err != nil {
				return
			}

			// Redirect
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			echoContext := e.NewContext(req, rec)
			echoContext.SetPath("/:id")
			echoContext.SetParamNames("id")
			echoContext.SetParamValues(shortenOut.Result)

			err = h.HandleShorten(echoContext)
			if err != nil {
				return
			}
		}
	})
}

func BenchmarkHandler_HandleUserUrls(b *testing.B) {
	b.Run("Endpoint: GET /api/user/urls", func(b *testing.B) {
		e := echo.New()
		h := NewHandler(ctx, service.NewService(storage.GetStorage(ctx)))

		for i := 0; i < b.N; i++ {
			// Post
			client := resty.New()

			for i := 0; i < 10; i++ {
				_, err := client.R().
					SetHeader(echo.HeaderContentType, echo.MIMETextPlain).
					SetBody([]byte(gofakeit.URL())).
					Post(config.GetConfig().BaseURL + "/")
				if err != nil {
					continue
				}
			}

			// Get
			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			rec := httptest.NewRecorder()
			echoContext := e.NewContext(req, rec)

			err := h.HandleShorten(echoContext)
			if err != nil {
				return
			}
		}
	})
}

func BenchmarkHandler_HandleUserUrlsDelete(b *testing.B) {
	b.Run("Endpoint: DELETE /api/user/urls", func(b *testing.B) {
		e := echo.New()
		h := NewHandler(ctx, service.NewService(storage.GetStorage(ctx)))

		for i := 0; i < b.N; i++ {
			// Post
			client := resty.New()

			for i := 0; i < 10; i++ {
				_, err := client.R().
					SetHeader(echo.HeaderContentType, echo.MIMETextPlain).
					SetBody([]byte(gofakeit.URL())).
					Post(config.GetConfig().BaseURL + "/")
				if err != nil {
					continue
				}
			}

			// Get
			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", nil)
			rec := httptest.NewRecorder()
			echoContext := e.NewContext(req, rec)

			err := h.HandleShorten(echoContext)
			if err != nil {
				return
			}
		}
	})
}
