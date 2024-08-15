package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
)

// HandleSave should handle empty or invalid request bodies
func TestHandleSaveEmptyOrInvalidBody(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	servAPI := service.NewService(storage.GetStorage())
	h := handler.NewHandler(ctx, servAPI)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.HandleSave(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Body is empty or invalid", rec.Body.String())
}

// Middleware logger should log incoming HTTP requests
func TestMiddlewareLogger(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.Background()

	logger.GetLogger()

	middleware := logger.MiddlewareLogger(ctx)
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
