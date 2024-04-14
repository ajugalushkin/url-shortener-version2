package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const create string = `CREATE TABLE IF NOT EXISTS shorten_urls (
    "short_url" VARCHAR(20) NOT NULL PRIMARY KEY,
    "correlation_id" VARCHAR(250) NOT NULL DEFAULT '',
    "original_url" VARCHAR(250) NOT NULL DEFAULT ''
) `

type PGStorage struct {
	db  *sql.DB
	ctx context.Context
}

func NewPGStorage(ctx context.Context) (*PGStorage, error) {
	flags := config.FlagsFromContext(ctx)
	db, err := sql.Open("pgx", flags.DataBaseDsn)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return &PGStorage{db: db, ctx: ctx}, nil
}

func (s *PGStorage) Put(shortening dto.Shortening) (*dto.Shortening, error) {
	_, err := s.db.ExecContext(
		s.ctx,
		"INSERT INTO shorten_urls (short_url,correlation_id,original_url) VALUES ($1,$2,$3)",
		shortening.ShortURL,
		shortening.CorrelationId,
		shortening.OriginalURL,
	)
	if err != nil {
		return nil, err
	}

	return &shortening, nil
}

func (s *PGStorage) Get(id string) (*dto.Shortening, error) {
	row := s.db.QueryRowContext(s.ctx, "SELECT * FROM shorten_urls WHERE short_url = ?", id)

	urlData := dto.Shortening{}
	if err := row.Scan(&urlData.ShortURL, &urlData.OriginalURL); err == sql.ErrNoRows {
		return &urlData, err
	}

	if urlData == (dto.Shortening{}) {
		return &urlData, errors.New("URL not found in DataBase")
	}

	return &urlData, nil
}
