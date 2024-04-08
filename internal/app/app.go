package app

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/ajugalushkin/url-shortener-version2/docs"
	"github.com/ajugalushkin/url-shortener-version2/internal/compress"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Run(ctx context.Context) error {
	flags := config.ConfigFromContext(ctx)

	log, err := logger.Initialize(flags.FlagLogLevel)
	if err != nil {
		return err
	}

	ctx = logger.ContextWithLogger(ctx, log)

	server := echo.New()
	serviceAPI := service.NewService(storage.NewStorage(ctx))
	newHandler := handler.NewHandler(ctx, serviceAPI)

	server.Use(logger.MiddlewareLogger(ctx))

	server.Use(compress.GzipWithConfig(compress.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	server.POST("/api/shorten", newHandler.HandleShorten)
	server.POST("/", newHandler.HandleSave)
	server.GET("/:id", newHandler.HandleRedirect)

	//Swagger
	server.GET("/swagger/*", echoSwagger.WrapHandler)

	fmt.Println("Running server on", flags.RunAddr)
	err = server.Start(flags.RunAddr)
	if err != nil {
		return err
	}

	return nil
}
