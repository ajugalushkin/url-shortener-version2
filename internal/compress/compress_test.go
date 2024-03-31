package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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

		context := server.NewContext(request, recorder)
		storageAPI := storage.NewInMemory()
		_, err = storageAPI.Put(model.Shortening{
			Key: "rIHY5pi",
			URL: "http://localhost:8080/rIHY5pi",
		})
		if assert.NoError(t, err) {
			middlewareGzip := GzipMiddleware(handler.NewHandler(service.NewService(storageAPI), &config.Config{}).HandleRedirect)

			// Assertions
			if assert.NoError(t, middlewareGzip(context)) {
				assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
			}
		}
	})

	//t.Run("accepts_gzip", func(t *testing.T) {
	//	buf := bytes.NewBufferString(requestBody)
	//
	//	r := httptest.NewRequest("POST", srv.URL, buf)
	//	r.RequestURI = ""
	//	r.Header.Set("Accept-Encoding", "gzip")
	//
	//	server := echo.New()
	//
	//	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/rIHY5pi", buf)
	//	request.RequestURI = ""
	//	request.Header.Set("Content-Encoding", "gzip")
	//	request.Header.Set("Accept-Encoding", "")
	//
	//	recorder := httptest.NewRecorder()
	//
	//	context := server.NewContext(request, recorder)
	//	storageAPI := storage.NewInMemory()
	//	_, err = storageAPI.Put(model.Shortening{
	//		Key: "rIHY5pi",
	//		URL: "http://localhost:8080/rIHY5pi",
	//	})
	//	if assert.NoError(t, err) {
	//		middlewareGzip := GzipMiddleware(handler.NewHandler(service.NewService(storageAPI), &config.Config{}).HandleRedirect)
	//
	//		// Assertions
	//		if assert.NoError(t, middlewareGzip(context)) {
	//			assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	//		}
	//	}
	//})
}
