package save

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostHandler(t *testing.T) {
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
			c := server.NewContext(req, rec)

			handler := New(service.NewService(storage.NewInMemory()))

			// Assertions
			if assert.NoError(t, handler(c)) {
				assert.Equal(t, test.want.code, rec.Code)
				assert.Equal(t, test.want.contentType, rec.Header().Get(echo.HeaderContentType))
			}
		})
	}
}
