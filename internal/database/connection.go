package database

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
)

// NewConnection функция для получения соединения с базой данных
func NewConnection(driver, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		logger.GetLogger().Error("failed to create a database connection", zap.Error(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.GetLogger().Error("failed to ping the database", zap.Error(err))
		return nil, err
	}

	Migrate(dsn)

	return db, nil
}
