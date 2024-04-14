package storage

import (
	"context"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
)

func GetStorage(ctx context.Context) service.PutGetter {
	flags := config.FlagsFromContext(ctx)
	if flags.DataBaseDsn != "" {
		storage, err := NewPGStorage(ctx)
		if err != nil {
			log := logger.LogFromContext(ctx)
			log.Error(err.Error())
			return nil
		}
		return storage
	} else if flags.FileStoragePath != "" {
		return NewStorage(ctx)
	} else {
		return NewInMemory()
	}
}
