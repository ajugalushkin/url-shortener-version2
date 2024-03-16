package app

import (
	"fmt"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/handlers"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config) error {
	server := echo.New()

	serviceAPI := service.NewService(storage.NewInMemory())

	handler := save.NewHandler(serviceAPI, cfg)

	server.POST("/", handler.HandleSave)
	server.GET("/:id", handler.HandleRedirect)

	fmt.Println("Running server on", cfg.FlagRunAddr)
	err := server.Start(cfg.FlagRunAddr)
	if err != nil {
		return err
	}

	return nil
}
