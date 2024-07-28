package compress

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
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

// Initializes compressWriter with a valid http.ResponseWriter
func TestNewCompressWriterWithValidResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	cw := newCompressWriter(w)

	if cw.w != w {
		t.Errorf("Expected http.ResponseWriter to be %v, got %v", w, cw.w)
	}

	if cw.zw == nil {
		t.Error("Expected gzip.Writer to be initialized, got nil")
	}
}

// Successfully creates a compressReader with a valid gzip stream
func TestNewCompressReaderWithValidGzipStream(t *testing.T) {
	// Create a buffer with gzip compressed data
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write([]byte("test data"))
	if err != nil {
		t.Fatalf("Failed to write gzip data: %v", err)
	}
	zw.Close()

	// Create a new compressReader
	r := io.NopCloser(&buf)
	cr, err := newCompressReader(r)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Read from compressReader and verify the output
	decompressedData := make([]byte, 100)
	n, err := cr.Read(decompressedData)
	if err != nil && err != io.EOF {
		t.Fatalf("Expected no error, got %v", err)
	}

	if string(decompressedData[:n]) != "test data" {
		t.Fatalf("Expected 'test data', got %s", string(decompressedData[:n]))
	}
}

// Handles invalid gzip data gracefully
func TestNewCompressReaderWithInvalidGzipData(t *testing.T) {
	// Create a buffer with invalid gzip data
	invalidData := []byte("invalid gzip data")
	r := io.NopCloser(bytes.NewReader(invalidData))

	// Attempt to create a new compressReader
	cr, err := newCompressReader(r)
	if cr != nil {
		t.Fatalf("Expected nil compressReader, got %v", cr)
	}

	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}

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

// Middleware skips compression when Skipper returns true
func TestMiddlewareSkipsCompressionWhenSkipperReturnsTrue(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := GzipConfig{
		Skipper: func(c echo.Context) bool {
			return true
		},
	}

	h := GzipWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotContains(t, rec.Header().Get("Content-Encoding"), "gzip")
	}
}

// Middleware handles nil Skipper by using default Skipper
func TestMiddlewareHandlesNilSkipperByUsingDefaultSkipper(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := GzipConfig{
		Skipper: nil,
	}

	h := GzipWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))
	}
}
