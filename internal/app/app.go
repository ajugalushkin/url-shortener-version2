package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	_ "github.com/ajugalushkin/url-shortener-version2/api"
	"github.com/ajugalushkin/url-shortener-version2/config"
	pb "github.com/ajugalushkin/url-shortener-version2/gen/url_shortener/v1"
	"github.com/ajugalushkin/url-shortener-version2/internal/compress"
	"github.com/ajugalushkin/url-shortener-version2/internal/handler"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"github.com/ajugalushkin/url-shortener-version2/pkg/ydx/url-shortener/v1"

	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
)

// Run является основным местом запуска сервиса.
// В методе происходит инициализация контекста, логгера и
// происходит привязка обработчиков к запросам.
func Run() error {
	mainCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer stop()

	serverHTTP := echo.New()
	setRouting(mainCtx, serverHTTP)

	group, groupCtx := errgroup.WithContext(mainCtx)
	group.Go(func() error {
		logger.GetLogger().Info("Running HTTP server", zap.String("address", config.GetConfig().ServerAddress))

		var err error
		if !config.GetConfig().EnableHTTPS {
			err = serverHTTP.Start(config.GetConfig().ServerAddress)
		} else {
			err = serverHTTP.StartAutoTLS(config.GetConfig().ServerAddress)
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.GetLogger().Fatal("shutting down the HTTP", zap.Error(err))
		}
		return err
	})
	group.Go(func() error {
		<-groupCtx.Done()
		if err := serverHTTP.Shutdown(context.Background()); err != nil {
			logger.GetLogger().Fatal(err.Error(), zap.String("address", config.GetConfig().ServerAddress))
			return err
		}
		return nil
	})

	serverGRPC := grpc.NewServer()
	newHandler := v1.NewHandler(mainCtx, service.NewService(storage.GetStorage()))
	pb.RegisterURLShortenerV1ServiceServer(serverGRPC, newHandler)
	group.Go(func() error {
		listen, err := net.Listen("tcp", config.GetConfig().ServerAddressGrpc)
		if err != nil {
			logger.GetLogger().Fatal("failed to listen", zap.Error(err))
			return err
		}

		logger.GetLogger().Debug("Running GRPC server", zap.String("address", config.GetConfig().ServerAddressGrpc))
		if err := serverGRPC.Serve(listen); err != nil {
			logger.GetLogger().Fatal("GRPC failed to serve", zap.Error(err))
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-groupCtx.Done()
		serverGRPC.Stop()
		logger.GetLogger().Debug("GRPC server has stopped")
		return nil
	})

	return group.Wait()
}

func setRouting(ctx context.Context, server *echo.Echo) {
	newHandler := handler.NewHandler(ctx, service.NewService(storage.GetStorage()))

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
	server.GET("/api/internal/stats", newHandler.FilterIP(newHandler.HandleStats))

	server.DELETE("/api/user/urls", newHandler.Authorized(newHandler.HandleUserUrlsDelete))

	//Swagger
	server.GET("/api/*", echoSwagger.WrapHandler)

	// Регистрация pprof-обработчиков
	server.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))

}
