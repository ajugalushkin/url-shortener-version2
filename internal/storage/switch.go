package storage

import (
	"context"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
)

func GetStorage(ctx context.Context) service.PutGetter {
	flags := config.ConfigFromContext(ctx)
	if flags.DataBaseDsn != "" {
		return NewPGStorage(ctx)
	} else if flags.FileStoragePath != "" {
		return NewStorage(ctx)
	} else {
		return NewInMemory()
	}
}
