package grpchandler

import (
	"context"
	"testing"

	"github.com/ajugalushkin/url-shortener-version2/internal/service"
)

// Initializes URLSServer with given context and service API
func TestNewHandlerInitialization(t *testing.T) {
	ctx := context.Background()
	servAPI := &service.Service{}

	handler := NewHandler(ctx, servAPI)

	if handler.ctx != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.ctx)
	}

	if handler.servAPI != servAPI {
		t.Errorf("expected service API %v, got %v", servAPI, handler.servAPI)
	}

	if handler.cache == nil {
		t.Error("expected cache to be initialized, got nil")
	}
}

// Handles nil context input gracefully
func TestNewHandlerNilContext(t *testing.T) {
	servAPI := &service.Service{}

	handler := NewHandler(nil, servAPI)

	if handler.ctx != nil {
		t.Errorf("expected context to be nil, got %v", handler.ctx)
	}

	if handler.servAPI != servAPI {
		t.Errorf("expected service API %v, got %v", servAPI, handler.servAPI)
	}

	if handler.cache == nil {
		t.Error("expected cache to be initialized, got nil")
	}
}
