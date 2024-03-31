package app

import (
	"fmt"
	"github.com/ajugalushkin/url-shortener-version2/internal/compress"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config) error {
	server := echo.New()

	serviceAPI := service.NewService(storage.NewStorage(cfg))
	newHandler := handler.NewHandler(serviceAPI, cfg)

	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		return err
	}

	server.Use(logger.RequestLogger)
	server.Use(compress.GzipMiddleware)

	server.POST("/api/shorten", newHandler.HandleSave)
	server.POST("/", newHandler.HandleSave)
	server.GET("/:id", newHandler.HandleRedirect)

	fmt.Println("Running server on", cfg.RunAddr)
	err := server.Start(cfg.RunAddr)
	if err != nil {
		return err
	}

	return nil
}
