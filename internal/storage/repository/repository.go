package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/ajugalushkin/url-shortener-version2/internal/database"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewRepository(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) Put(ctx context.Context, shortening dto.Shortening) (*dto.Shortening, error) {
	var (
		result sql.Result
		err    error
	)

	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Insert("shorten_urls").
			Columns("short_url", "correlation_id", "original_url").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		sb = sb.Values(
			shortening.ShortURL,
			shortening.CorrelationID,
			shortening.OriginalURL,
		)

		result, err = sb.ExecContext(ctx)
		return err
	})

	if err != nil {
		return nil, errors.Wrap(err, "repository.Put")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Put")
	}
	shortening.ShortURL = strconv.FormatInt(id, 10)

	return &shortening, nil
}

func (r *Repo) Get(ctx context.Context, shortURL string) (*dto.Shortening, error) {
	var shorteningList []dto.Shortening

	err := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.Select("short_url", "correlation_id", "original_url").
			From("shorten_urls").
			PlaceholderFormat(squirrel.Dollar).
			Where(squirrel.Eq{"short_url": []string{shortURL}}).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &shorteningList, query, args...)
	})

	if err != nil {
		return nil, errors.Wrap(err, "repository.Get")
	}

	if len(shorteningList) == 0 {
		return nil, errors.Wrap(sql.ErrNoRows, "repository.Get")
	}

	shortening := shorteningList[0]

	return &shortening, nil
}

func (r *Repo) PutList(ctx context.Context, list dto.ShorteningList) error {
	var err error
	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Insert("shorten_urls").
			Columns("short_url", "correlation_id", "original_url").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		for _, shortening := range list {
			sb = sb.Values(
				shortening.ShortURL,
				shortening.CorrelationID,
				shortening.OriginalURL,
			)
		}

		_, err = sb.ExecContext(ctx)
		return err
	})

	if err != nil {
		return errors.Wrap(err, "repository.CreateGoods")
	}

	return nil
}
