package handler

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var cfg = config.Config{
	RunAddr: "localhost:8080",
	BaseURL: "http://localhost:8080",
}

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
			context := server.NewContext(req, rec)

			storageAPI := storage.NewInMemory()
			_, err := storageAPI.Put(model.Shortening{
				Key: test.request.key,
				URL: test.want.response,
			})
			if assert.NoError(t, err) {
				handler := NewHandler(service.NewService(storageAPI), &cfg)

				// Assertions
				if assert.NoError(t, handler.HandleRedirect(context)) {
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
		{
			name: "Test JSON Ok",
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			server := echo.New()
			req := httptest.NewRequest(test.request.method, "/", strings.NewReader(test.request.body))
			req.Header.Set(echo.HeaderContentType, test.request.contentType)
			rec := httptest.NewRecorder()
			context := server.NewContext(req, rec)

			handler := NewHandler(service.NewService(storage.NewInMemory()), &cfg)

			// Assertions
			if assert.NoError(t, handler.HandleSave(context)) {
				assert.Equal(t, test.want.code, rec.Code)
				assert.Equal(t, test.want.contentType, rec.Header().Get(echo.HeaderContentType))
			}
		})
	}
}

//func TestHandler_HandleShorten(t *testing.T) {
//	type request struct {
//		method      string
//		body        string
//		contentType string
//	}
//	type want struct {
//		code        int
//		contentType string
//	}
//	tests := []struct {
//		name    string
//		request request
//		want    want
//	}{
//		{
//			name: "Test StatusCreated",
//			request: request{
//				method:      http.MethodPost,
//				body:        "{\n  \"url\": \"https://practicum.yandex.ru\"\n}",
//				contentType: echo.MIMEApplicationJSON,
//			},
//			want: want{
//				code:        http.StatusCreated,
//				contentType: echo.MIMEApplicationJSON,
//			},
//		},
//		{
//			name: "Test Wrong Request",
//			request: request{
//				method:      http.MethodGet,
//				body:        "{\n  \"url\": \"https://practicum.yandex.ru\"\n}",
//				contentType: echo.MIMEApplicationJSON,
//			},
//			want: want{
//				code:        http.StatusBadRequest,
//				contentType: echo.MIMETextPlainCharsetUTF8,
//			},
//		},
//		{
//			name: "Test Empty URL",
//			request: request{
//				method:      http.MethodGet,
//				body:        "",
//				contentType: echo.MIMEApplicationJSON,
//			},
//			want: want{
//				code:        http.StatusBadRequest,
//				contentType: echo.MIMETextPlainCharsetUTF8,
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			// Setup
//			server := echo.New()
//			req := httptest.NewRequest(test.request.method, "/api/shorten", strings.NewReader(test.request.body))
//			req.Header.Set(echo.HeaderContentType, test.request.contentType)
//			rec := httptest.NewRecorder()
//			context := server.NewContext(req, rec)
//
//			handler := NewHandler(service.NewService(storage.NewInMemory()), &cfg)
//
//			// Assertions
//			if assert.NoError(t, handler.HandleShorten(context)) {
//				assert.Equal(t, test.want.code, rec.Code)
//				assert.Equal(t, test.want.contentType, rec.Header().Get(echo.HeaderContentType))
//			}
//		})
//	}
//}
