package redirect

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	type want struct {
		code     int
		request  string
		method   string
		response string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Test Bad Request 1",
			want: want{
				code:    http.StatusBadRequest,
				request: "https://practicum.yandex.ru/",
				method:  http.MethodGet,
			},
		},
		{
			name: "Test Bad Request 2",
			want: want{
				code:    http.StatusBadRequest,
				request: "https://practicum.yandex.ru/",
				method:  http.MethodPost,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.want.method, test.want.request, nil)
			request.Header.Set("Content-Type", "text/plain")
			writer := httptest.NewRecorder()

			serviceAPI := service.NewService(storage.NewInMemory())

			handler := New(serviceAPI)
			handler.ServeHTTP(writer, request)

			result := writer.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.code, result.StatusCode)
			assert.Equal(t, test.want.response, result.Header.Get("Location"))
		})
	}

	tests = []struct {
		name string
		want want
	}{
		{
			name: "Test Bad Request 1",
			want: want{
				code:    http.StatusBadRequest,
				request: "https://practicum.yandex.ru/",
				method:  http.MethodGet,
			},
		},
		{
			name: "Test Bad Request 2",
			want: want{
				code:    http.StatusBadRequest,
				request: "https://practicum.yandex.ru/",
				method:  http.MethodPost,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.want.method, test.want.request, nil)
			request.Header.Set("Content-Type", "text/plain")
			writer := httptest.NewRecorder()

			handler := New(service.NewService(storage.NewInMemory()))
			handler.ServeHTTP(writer, request)

			result := writer.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.code, result.StatusCode)
			assert.Equal(t, test.want.response, result.Header.Get("Location"))
		})
	}
}
