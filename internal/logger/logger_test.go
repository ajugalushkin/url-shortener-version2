package logger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	//"go.uber.org/zap"
)

// Logs incoming HTTP request details
func TestMiddlewareLoggerLogsRequestDetails(t *testing.T) {
	// Initialize the logger
	logger := GetLogger()

	// Create a new echo instance
	e := echo.New()

	// Create a new request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a middleware function
	middleware := MiddlewareLogger(context.Background())

	// Create a handler function
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Wrap the handler with the middleware
	wrappedHandler := middleware(handler)

	// Call the wrapped handler
	if assert.NoError(t, wrappedHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
	}

	// Check the logs
	logs := logger.Check(zap.DebugLevel, "got incoming HTTP request")
	assert.NotNil(t, logs)
}

// Handles nil context
func TestMiddlewareLoggerHandlesNilContext(t *testing.T) {
	// Initialize the logger
	logger := GetLogger()

	// Create a new echo instance
	e := echo.New()

	// Create a new request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a middleware function with nil context
	middleware := MiddlewareLogger(nil)

	// Create a handler function
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Wrap the handler with the middleware
	wrappedHandler := middleware(handler)

	// Call the wrapped handler
	if assert.NoError(t, wrappedHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
	}

	// Check the logs
	logs := logger.Check(zap.DebugLevel, "got incoming HTTP request")
	assert.NotNil(t, logs)
}
