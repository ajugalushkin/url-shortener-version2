package repository

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/ajugalushkin/url-shortener-version2/internal/database"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	userErr "github.com/ajugalushkin/url-shortener-version2/internal/errors"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func NewRepository(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) Put(ctx context.Context, shorteningInput dto.Shortening) (*dto.Shortening, error) {
	var err error
	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Insert("shorten_urls").
			Columns("short_url", "correlation_id", "original_url").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		sb = sb.Values(
			shorteningInput.ShortURL,
			shorteningInput.CorrelationID,
			shorteningInput.OriginalURL,
		)

		_, err = sb.ExecContext(ctx)
		return err
	})

	if err != nil {
		if pgErr, ok := errors.Unwrap(errors.Unwrap(err)).(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			shortening, _ := r.GetByURL(ctx, shorteningInput.OriginalURL)
			if shortening.OriginalURL != "" {
				return shortening, errors.Wrapf(userErr.ErrorDuplicateURL, "%s %s", userErr.ErrorDuplicateURL, shortening.OriginalURL)
			}
		}
		return nil, errors.Wrap(err, "repository.Put")
	}
	return &shorteningInput, nil
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
	log := logger.LogFromContext(ctx)
	if err != nil {
		log.Info("repository.Get", zap.Error(err))
		return nil, errors.Wrap(err, "repository.Get")
	}

	if len(shorteningList) == 0 {
		log.Info("repository.Get", zap.Error(sql.ErrNoRows))
		return nil, errors.Wrap(sql.ErrNoRows, "repository.Get")
	}

	shortening := shorteningList[0]

	log.Info("repository.Get OK", zap.String("Original URL", shortening.OriginalURL))

	return &shortening, nil
}

func (r *Repo) GetByURL(ctx context.Context, originURL string) (*dto.Shortening, error) {
	var shorteningList []dto.Shortening

	err := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.Select("short_url", "correlation_id", "original_url").
			From("shorten_urls").
			PlaceholderFormat(squirrel.Dollar).
			Where(squirrel.Eq{"original_url": []string{originURL}}).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &shorteningList, query, args...)
	})
	log := logger.LogFromContext(ctx)
	if err != nil {
		log.Info("repository.Get", zap.Error(err))
		return nil, errors.Wrap(err, "repository.Get")
	}

	if len(shorteningList) == 0 {
		log.Info("repository.Get", zap.Error(sql.ErrNoRows))
		return nil, errors.Wrap(sql.ErrNoRows, "repository.Get")
	}

	shortening := shorteningList[0]

	log.Info("repository.Get OK", zap.String("Original URL", shortening.OriginalURL))

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
		return errors.Wrap(err, "repository.PutList")
	}

	return nil
}
