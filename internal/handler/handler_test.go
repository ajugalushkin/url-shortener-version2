package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var newConfig = config.Config{
	RunAddr: "localhost:8080",
	BaseURL: "http://localhost:8080",
}

var ctx = config.ContextWithFlags(context.Background(), &newConfig)

func TestHandler_HandleRedirect(t *testing.T) {
	type request struct {
		method      string
		key         string
		URL         string
		contentType string
	}
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Test OK",
			request: request{
				method:      http.MethodGet,
				key:         "rIHY5pi",
				URL:         "http://localhost:8080/rIHY5pi",
				contentType: echo.MIMETextPlain,
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "https://practicum.yandex.ru/",
			},
		},
		{
			name: "Test Bad Request 1",
			request: request{
				method:      http.MethodGet,
				URL:         "http://localhost:8080/rIHY5pi",
				contentType: echo.MIMETextPlain,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Test Bad Request 2",
			request: request{
				URL:    "http://localhost:8080/rIHY5pi",
				method: http.MethodPost,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			server := echo.New()
			req := httptest.NewRequest(test.request.method, test.request.URL, nil)
			req.Header.Set(echo.HeaderContentType, test.request.contentType)
			rec := httptest.NewRecorder()
			echoCtx := server.NewContext(req, rec)

			storageAPI := storage.NewInMemory()
			_, err := storageAPI.Put(dto.Shortening{
				ShortURL:    test.request.key,
				OriginalURL: test.want.response,
			})
			if assert.NoError(t, err) {
				handler := NewHandler(ctx, service.NewService(storageAPI))

				// Assertions
				if assert.NoError(t, handler.HandleRedirect(echoCtx)) {
					assert.Equal(t, test.want.code, rec.Code)
					assert.Equal(t, test.want.response, rec.Header().Get(echo.HeaderLocation))
				}
			}
		})
	}
}

func TestHandler_HandleSave(t *testing.T) {
	type request struct {
		method      string
		body        string
		contentType string
	}
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Test StatusCreated",
			request: request{
				method:      http.MethodPost,
				body:        "https://practicum.yandex.ru/",
				contentType: echo.MIMETextPlain,
			},
			want: want{
				code:        http.StatusCreated,
				contentType: echo.MIMETextPlain,
			},
		},
		{
			name: "Test Wrong Request",
			request: request{
				method:      http.MethodGet,
				body:        "https://practicum.yandex.ru/",
				contentType: echo.MIMETextPlain,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: echo.MIMETextPlainCharsetUTF8,
			},
		},
		{
			name: "Test Empty URL",
			request: request{
				method:      http.MethodGet,
				body:        "",
				contentType: echo.MIMETextPlain,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: echo.MIMETextPlainCharsetUTF8,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			server := echo.New()
			req := httptest.NewRequest(test.request.method, "/", strings.NewReader(test.request.body))
			req.Header.Set(echo.HeaderContentType, test.request.contentType)
			rec := httptest.NewRecorder()
			echoCtx := server.NewContext(req, rec)

			handler := NewHandler(ctx, service.NewService(storage.NewInMemory()))

			// Assertions
			if assert.NoError(t, handler.HandleSave(echoCtx)) {
				assert.Equal(t, test.want.code, rec.Code)
				assert.Equal(t, test.want.contentType, rec.Header().Get(echo.HeaderContentType))
			}
		})
	}
}

func TestHandler_HandleShorten(t *testing.T) {
	type request struct {
		method      string
		body        string
		contentType string
	}
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Test StatusCreated",
			request: request{
				method:      http.MethodPost,
				body:        "{\n  \"url\": \"https://practicum.yandex.ru\"\n}",
				contentType: echo.MIMEApplicationJSON,
			},
			want: want{
				code:        http.StatusCreated,
				contentType: echo.MIMEApplicationJSON,
			},
		},
		{
			name: "Test Wrong Request",
			request: request{
				method:      http.MethodGet,
				body:        "{\n  \"url\": \"https://practicum.yandex.ru\"\n}",
				contentType: echo.MIMEApplicationJSON,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: echo.MIMETextPlainCharsetUTF8,
			},
		},
		{
			name: "Test Empty URL",
			request: request{
				method:      http.MethodGet,
				body:        "",
				contentType: echo.MIMEApplicationJSON,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: echo.MIMETextPlainCharsetUTF8,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			server := echo.New()
			req := httptest.NewRequest(test.request.method, "/api/shorten", strings.NewReader(test.request.body))
			req.Header.Set(echo.HeaderContentType, test.request.contentType)
			rec := httptest.NewRecorder()
			echoCtx := server.NewContext(req, rec)

			handler := NewHandler(ctx, service.NewService(storage.NewInMemory()))

			// Assertions
			if assert.NoError(t, handler.HandleShorten(echoCtx)) {
				assert.Equal(t, test.want.code, rec.Code)
				assert.Equal(t, test.want.contentType, rec.Header().Get(echo.HeaderContentType))
			}
		})
	}
}

func TestHandler_HandleShortenBatch(t *testing.T) {
	type fields struct {
		ctx     context.Context
		servAPI *service.Service
	}
	tests := []struct {
		name             string
		fields           fields
		inputContentType string
		inputMethod      string
		inputBody        string
		expectedHeader   string
		expectedCode     int
	}{
		{
			name: "Test ОК",
			fields: fields{
				ctx:     ctx,
				servAPI: service.NewService(storage.NewInMemory())},
			inputContentType: echo.MIMEApplicationJSON,
			inputMethod:      http.MethodPost,
			inputBody:        "[\n    {\n        \"correlation_id\": \"1\",\n        \"original_url\": \"https://vk.com/ajugalushkin\"\n    }\n]",
			expectedHeader:   echo.MIMEApplicationJSON,
			expectedCode:     http.StatusCreated,
		},
		{
			name: "Test Bad Request Type",
			fields: fields{
				ctx:     ctx,
				servAPI: service.NewService(storage.NewInMemory())},
			inputContentType: echo.MIMEApplicationJSON,
			inputMethod:      http.MethodGet,
			inputBody:        "[\n    {\n        \"correlation_id\": \"1\",\n        \"original_url\": \"https://vk.com/ajugalushkin\"\n    }\n]",
			expectedHeader:   echo.MIMETextPlainCharsetUTF8,
			expectedCode:     http.StatusBadRequest,
		},
		{
			name: "Test Bad Content Type",
			fields: fields{
				ctx:     ctx,
				servAPI: service.NewService(storage.NewInMemory())},
			inputContentType: echo.MIMETextPlain,
			inputMethod:      http.MethodPost,
			inputBody:        "[\n    {\n        \"correlation_id\": \"1\",\n        \"original_url\": \"https://vk.com/ajugalushkin\"\n    }\n]",
			expectedHeader:   echo.MIMETextPlainCharsetUTF8,
			expectedCode:     http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			req := httptest.NewRequest(
				test.inputMethod,
				"/api/shorten/batch",
				strings.NewReader(test.inputBody),
			)
			req.Header.Set(echo.HeaderContentType, test.inputContentType)
			rec := httptest.NewRecorder()

			echoCtx := echo.New().NewContext(req, rec)

			handler := Handler{ctx: test.fields.ctx, servAPI: test.fields.servAPI}

			// Assertions
			if assert.NoError(t, handler.HandleShortenBatch(echoCtx)) {
				assert.Equal(t, test.expectedHeader, rec.Header().Get(echo.HeaderContentType))
				assert.Equal(t, test.expectedCode, rec.Code)
			}
		})
	}
}
