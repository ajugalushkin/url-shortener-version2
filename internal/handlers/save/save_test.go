package save

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
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
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name: "Test Wrong Request",
			request: request{
				method:      http.MethodGet,
				body:        "https://practicum.yandex.ru/",
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, "/",
				strings.NewReader(test.request.body))
			request.Header.Set("Content-Type", test.request.contentType)

			writer := httptest.NewRecorder()

			handler := New(service.NewService(storage.NewInMemory()))
			handler.ServeHTTP(writer, request)

			result := writer.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.code, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
