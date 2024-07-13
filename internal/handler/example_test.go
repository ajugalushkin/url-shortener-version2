package handler

import (
	"fmt"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/inmemory"
)

func ExampleRedirectHandler() {
	server := echo.New()
	req := httptest.NewRequest(echo.GET, "http://localhost:8080/sg45fw", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	echoCtx := server.NewContext(req, rec)

	storageAPI := inmemory.NewInMemory()
	_, err := storageAPI.Put(ctx, dto.Shortening{
		ShortURL:    "sg45fw",
		OriginalURL: "http://test.ru",
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	handler := NewHandler(ctx, service.NewService(storageAPI))
	err = handler.HandleRedirect(echoCtx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
