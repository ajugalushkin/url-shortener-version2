package storage

import (
	"context"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/database"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/file"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/inmemory"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/repository"
)

func GetStorage(ctx context.Context) service.PutGetter {
	flags := config.FlagsFromContext(ctx)
	if flags.DataBaseDsn != "" {
		db, err := database.NewConnection("pgx", flags.DataBaseDsn)
		if err != nil {
			log := logger.LogFromContext(ctx)
			log.Error(err.Error())
			return nil
		}

		repo := repository.NewRepository(db)
		if err != nil {
			log := logger.LogFromContext(ctx)
			log.Error(err.Error())
			return nil
		}
		return repo
	} else if flags.FileStoragePath != "" {
		return file.NewStorage(flags.FileStoragePath)
	} else {
		return inmemory.NewInMemory()
	}
}
