package compress

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/inmemory"
)

func TestGzipMiddleware(t *testing.T) {
	requestBody := `https://practicum.yandex.ru/`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		server := echo.New()

		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/rIHY5pi", buf)
		request.RequestURI = ""
		request.Header.Set("Content-Encoding", "gzip")
		request.Header.Set("Accept-Encoding", "")

		recorder := httptest.NewRecorder()

		echoCtx := server.NewContext(request, recorder)
		storageAPI := inmemory.NewInMemory()
		ctx := context.Background()
		_, err = storageAPI.Put(ctx, dto.Shortening{
			ShortURL:    "rIHY5pi",
			OriginalURL: "http://localhost:8080/rIHY5pi",
		})
		if assert.NoError(t, err) {
			handlerGzip := Gzip()
			middlewareGzip := handlerGzip(handler.NewHandler(ctx, service.NewService(storageAPI)).HandleRedirect)

			// Assertions
			if assert.NoError(t, middlewareGzip(echoCtx)) {
				assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
			}
		}
	})
}
