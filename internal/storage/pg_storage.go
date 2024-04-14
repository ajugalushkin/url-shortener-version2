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
		shortening.CorrelationID,
		shortening.OriginalURL,
	)
	if err != nil {
		return nil, err
	}

	return &shortening, nil
}

func (s *PGStorage) PutList(list dto.ShorteningList) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for _, item := range list {
		res, errNotFound := s.Get(item.ShortURL)
		if errNotFound == nil && *res != (dto.Shortening{}) {
			continue
		}

		_, err := tx.ExecContext(s.ctx,
			"INSERT INTO shorten_urls (short_url,correlation_id,original_url) "+
				"VALUES ($1,$2,$3)", item.ShortURL, item.CorrelationID, item.OriginalURL)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return errRollback
			}
			return err
		}
	}
	return tx.Commit()
}

func (s *PGStorage) Get(id string) (*dto.Shortening, error) {
	var shortening dto.Shortening

	row := s.db.QueryRowContext(s.ctx, "SELECT * FROM shorten_urls WHERE short_url = $1", id)
	if row.Err() != nil {
		return &shortening, row.Err()
	}

	if err := row.Scan(&shortening.ShortURL, &shortening.CorrelationID, &shortening.OriginalURL); err == sql.ErrNoRows {
		return &shortening, err
	}

	if shortening == (dto.Shortening{}) {
		return &shortening, errors.New("URL not found in DataBase")
	}

	return &shortening, nil
}
