package logger

import (
	"context"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ajugalushkin/url-shortener-version2/config"
)

// переменные для генерации инстанции
var (
	logger *zap.Logger
	once   sync.Once
)

// GetLogger создает инстанцию логгера
func GetLogger() *zap.Logger {
	once.Do(
		func() {
			// инициализируем объект
			log, err := initialize(config.GetConfig().FlagLogLevel)
			if err != nil {
				logger = zap.L()
				return
			}
			logger = log
		})

	return logger
}

// Initialize функция инициализации
func initialize(level string) (*zap.Logger, error) {
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

// MiddlewareLogger функция middleware
func MiddlewareLogger(ctx context.Context) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		fn := func(context echo.Context) error {
			start := time.Now()

			if err := next(context); err != nil {
				context.Error(err)
			}

			duration := time.Since(start)

			//log := LogFromContext(ctx)
			logger.Debug("got incoming HTTP request",
				zap.String("method", context.Request().Method),
				zap.String("path", context.Request().URL.Path),
				zap.String("time", duration.String()),
			)
			return nil
		}
		return fn
	}
}
