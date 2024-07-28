package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
)

// Server starts successfully with HTTP
//func TestServerStartsSuccessfullyWithHTTP(t *testing.T) {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	go func() {
//		time.Sleep(1 * time.Second)
//		cancel()
//	}()
//
//	err := Run(ctx)
//	if err != nil {
//		t.Fatalf("Expected no error, got %v", err)
//	}
//}

// Server fails to start due to invalid server address
//func TestServerFailsToStartDueToInvalidServerAddress(t *testing.T) {
//	originalAddress := config.GetConfig().ServerAddress
//	config.GetConfig().ServerAddress = "invalid_address"
//	defer func() { config.GetConfig().ServerAddress = originalAddress }()
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	err := Run(ctx)
//	if err == nil {
//		t.Fatalf("Expected error due to invalid server address, got nil")
//	}
//}

// Middleware logger should log incoming HTTP requests
//func TestMiddlewareLogger(t *testing.T) {
//	e := echo.New()
//	ctx := context.Background()
//	req := httptest.NewRequest(http.MethodGet, "/", nil)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	mw := logger.MiddlewareLogger(ctx)
//	handler := mw(func(c echo.Context) error {
//		return c.String(http.StatusOK, "test")
//	})
//
//	err := handler(c)
//	assert.NoError(t, err)
//	assert.Equal(t, http.StatusOK, rec.Code)
//}

// HandleSave should handle empty or invalid request bodies
func TestHandleSaveEmptyOrInvalidBody(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	servAPI := service.NewService(storage.GetStorage(ctx))
	h := handler.NewHandler(ctx, servAPI)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.HandleSave(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Body is empty or invalid", rec.Body.String())
}
