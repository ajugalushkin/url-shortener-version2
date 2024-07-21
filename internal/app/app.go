package app

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	_ "github.com/ajugalushkin/url-shortener-version2/api"
	"github.com/ajugalushkin/url-shortener-version2/internal/compress"
	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
)

// Run является основным местом запуска сервиса.
// В методе происходит инициализация контекста, логгера и
// происходит привязка обработчиков к запросам.
func Run(ctx context.Context) error {
	flags := config.FlagsFromContext(ctx)

	log, err := logger.Initialize(flags.FlagLogLevel)
	if err != nil {
		return err
	}

	ctx = logger.ContextWithLogger(ctx, log)

	server := echo.New()
	setRouting(ctx, server)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer stop()

	go func() {
		log.Info("Running server", zap.String("address", flags.ServerAddress))

		var err error
		if !flags.EnableHTTPS {
			err = server.Start(flags.ServerAddress)
		} else {
			err = server.StartAutoTLS(flags.ServerAddress)
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("shutting down the server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err.Error(), zap.String("address", flags.ServerAddress))
	}

	return nil
}

func setRouting(ctx context.Context, server *echo.Echo) {
	newHandler := handler.NewHandler(ctx, service.NewService(storage.GetStorage(ctx)))

	//Middleware
	server.Use(logger.MiddlewareLogger(ctx))
	server.Use(compress.GzipWithConfig(compress.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "api") ||
				strings.Contains(c.Request().URL.Path, "debug")
		},
	}))

	//Handlers
	server.POST("/", newHandler.HandleSave)
	server.POST("/api/shorten", newHandler.HandleShorten)
	server.POST("/api/shorten/batch", newHandler.HandleShortenBatch)

	server.GET("/:id", newHandler.HandleRedirect)
	server.GET("/ping", newHandler.HandlePing)
	server.GET("/api/user/urls", newHandler.Authorized(newHandler.HandleUserUrls))

	server.DELETE("/api/user/urls", newHandler.Authorized(newHandler.HandleUserUrlsDelete))

	//Swagger
	server.GET("/api/*", echoSwagger.WrapHandler)

	// Регистрация pprof-обработчиков
	server.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))

}
