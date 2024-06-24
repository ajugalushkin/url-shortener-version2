package logger

import (
	"context"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ajugalushkin/url-shortener-version2/config"
)

type singleInstance struct {
	logger *zap.Logger
}

func (i singleInstance) GetLogger() *zap.Logger {
	return i.logger
}

var (
	singleton *singleInstance
	once      sync.Once
)

func GetSingleton(ctx context.Context) *singleInstance {
	once.Do(
		func() {
			logger, err := Initialize(config.FlagsFromContext(ctx).FlagLogLevel)
			if err != nil {
				logger = zap.L()
			}
			singleton = &singleInstance{logger: logger}
		})

	return singleton
}

func Initialize(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zl, nil
}

func MiddlewareLogger(ctx context.Context) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		fn := func(context echo.Context) error {
			start := time.Now()

			if err := next(context); err != nil {
				context.Error(err)
			}

			duration := time.Since(start)

			log := GetSingleton(ctx).logger
			log.Debug("got incoming HTTP request",
				zap.String("method", context.Request().Method),
				zap.String("path", context.Request().URL.Path),
				zap.String("time", duration.String()),
			)
			return nil
		}
		return fn
	}
}
