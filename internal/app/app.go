package app

import (
	"fmt"
	_ "github.com/ajugalushkin/url-shortener-version2/docs"
	"github.com/ajugalushkin/url-shortener-version2/internal/compress"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"strings"
)

func Run(cfg *config.Config) error {
	server := echo.New()

	serviceAPI := service.NewService(storage.NewStorage(cfg))
	newHandler := handler.NewHandler(serviceAPI, cfg)

	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		return err
	}

	server.Use(logger.RequestLogger)
	server.Use(compress.GzipWithConfig(compress.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
	}))

	server.POST("/api/shorten", newHandler.HandleShorten)
	server.POST("/", newHandler.HandleSave)
	server.GET("/:id", newHandler.HandleRedirect)

	//Swagger
	server.GET("/swagger/*", echoSwagger.WrapHandler)

	fmt.Println("Running server on", cfg.RunAddr)
	err := server.Start(cfg.RunAddr)
	if err != nil {
		return err
	}

	return nil
}
