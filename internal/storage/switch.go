package storage

import (
	"go.uber.org/zap"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/database"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/file"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/inmemory"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/repository"
)

// GetStorage функция определяет инстанцию хранилища в зависимости от конфигурации.
func GetStorage() service.PutGetter {
	flags := config.GetConfig()
	log := logger.GetLogger()
	if flags.DataBaseDsn != "" {
		db, err := database.NewConnection("pgx", flags.DataBaseDsn)
		if err != nil {
			log.Error("storage.GetStorage Error:", zap.Error(err))
			return nil
		}
		log.Info("storage.GetStorage Set PostgresSQL")
		return repository.NewRepository(db)
	} else if flags.FileStoragePath != "" {
		log.Info("storage.GetStorage Set File")
		return file.NewStorage(flags.FileStoragePath)
	} else {
		log.Info("storage.GetStorage Set In Memory")
		return inmemory.NewInMemory()
	}
}
