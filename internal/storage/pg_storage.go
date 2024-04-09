package storage

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const create string = `CREATE TABLE IF NOT EXISTS urls (
	key	TEXT NOT NULL PRIMARY key,
	url	TEXT NOT NULL 
)`

type PGStorage struct {
	db  *sql.DB
	ctx context.Context
}

func NewPGStorage(ctx context.Context) *PGStorage {
	flags := config.ConfigFromContext(ctx)
	db, err := sql.Open("pgx", flags.DataBaseDsn)
	if err != nil {
		return nil
	}
	if _, err := db.Exec(create); err != nil {
		return nil
	}
	return &PGStorage{
		db:  db,
		ctx: ctx,
	}
}

func (s *PGStorage) Put(shortening dto.Shortening) (*dto.Shortening, error) {
	res, err := s.db.ExecContext(s.ctx, "INSERT INTO urls (key,url) VALUES ($1,$2)",
		shortening.Key, shortening.URL)
	if err != nil {
		return nil, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}

	shortening.Key = strconv.FormatInt(id, 10)

	return &shortening, nil
}

func (s *PGStorage) Get(id string) (*dto.Shortening, error) {
	row := s.db.QueryRowContext(s.ctx, "SELECT * FROM urls WHERE KEY = $1", id)

	urlData := dto.Shortening{}
	if err := row.Scan(&urlData.Key, &urlData.URL); err == sql.ErrNoRows {
		return &urlData, err
	}
	return &urlData, nil
}
