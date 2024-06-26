package app

import (
	"context"
	_ "github.com/ajugalushkin/url-shortener-version2/docs"
	"github.com/ajugalushkin/url-shortener-version2/internal/compress"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

func Run(ctx context.Context) error {
	flags := config.FlagsFromContext(ctx)

	log, err := logger.Initialize(flags.FlagLogLevel)
	if err != nil {
		return err
	}

	ctx = logger.ContextWithLogger(ctx, log)

	server := echo.New()
	newHandler := handler.NewHandler(ctx, service.NewService(storage.GetStorage(ctx)))

	//Middleware
	server.Use(logger.MiddlewareLogger(ctx))
	server.Use(compress.GzipWithConfig(compress.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger") ||
				strings.Contains(c.Request().URL.Path, "debug")
		},
	}))

	//Handlers
	server.POST("/", newHandler.HandleSave)
	server.POST("/api/shorten", newHandler.HandleShorten)
	server.POST("/api/shorten/batch", newHandler.HandleShortenBatch)

	server.GET("/:id", newHandler.HandleRedirect)
	server.GET("/ping", newHandler.HandlePing)
	server.GET("/api/user/urls", newHandler.HandleUserUrls)

	server.DELETE("/api/user/urls", newHandler.HandleUserUrlsDelete)

	//Swagger
	server.GET("/swagger/*", echoSwagger.WrapHandler)

	// Регистрация pprof-обработчиков
	server.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))

	log.Info("Running server", zap.String("address", flags.RunAddr))
	err = server.Start(flags.RunAddr)
	if err != nil {
		return err
	}

	return nil
}
