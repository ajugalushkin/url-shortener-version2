package logger

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		start := time.Now()

		if err := next(context); err != nil {
			context.Error(err)
		}

		duration := time.Since(start)

		Log.Debug("got incoming HTTP request",
			zap.String("method", context.Request().Method),
			zap.String("path", context.Request().URL.Path),
			zap.String("time", duration.String()),
		)
		return nil
	}
}
