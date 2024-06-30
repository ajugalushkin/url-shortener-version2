package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// описание типов
type (
	// Skipper функция для пропуска сжатия для swagger.
	Skipper func(c echo.Context) bool

	// GzipConfig структура для пропуска сжатия для swagger.
	GzipConfig struct {
		Skipper Skipper
	}
)

// объявление переменных
var (
	// DefaultGzipConfig структура по умолчанию для отключения сжатия.
	DefaultGzipConfig = GzipConfig{
		Skipper: func(c echo.Context) bool {
			return false
		},
	}
)

// compressWriter структура компрессии.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// newCompressWriter конструктор.
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header реализация вызова Header.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write запись.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader запись заголовка.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressReader конструктор.
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read чтение.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрытие.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Gzip объявление функции middleware
func Gzip() echo.MiddlewareFunc {
	return GzipWithConfig(DefaultGzipConfig)
}

// GzipWithConfig функция middleware
func GzipWithConfig(config GzipConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultGzipConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			if config.Skipper(context) {
				return next(context)
			}

			ow := context.Response().Writer

			acceptEncoding := context.Request().Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				cw := newCompressWriter(context.Response().Writer)
				ow = cw
				defer cw.Close()
			}

			contentEncoding := context.Request().Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(context.Request().Body)
				if err != nil {
					context.Response().WriteHeader(http.StatusInternalServerError)
					return nil
				}

				context.Request().Body = cr
				defer cr.Close()
			}

			context.Response().Writer = ow
			if err := next(context); err != nil {
				context.Error(err)
			}
			return nil
		}
	}
}
