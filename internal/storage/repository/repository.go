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
	"strconv"
	"sync"
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
			Columns("short_url", "correlation_id", "original_url", "user_id").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		sb = sb.Values(
			shorteningInput.ShortURL,
			shorteningInput.CorrelationID,
			shorteningInput.OriginalURL,
			shorteningInput.UserID,
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
		sb := squirrel.Select("short_url", "correlation_id", "original_url", "user_id", "is_deleted").
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
		log.Debug("repository.Get", zap.Error(err))
		return nil, errors.Wrap(err, "repository.Get")
	}

	if len(shorteningList) == 0 {
		log.Debug("repository.Get", zap.Error(sql.ErrNoRows))
		return nil, errors.Wrap(sql.ErrNoRows, "repository.Get")
	}

	shortening := shorteningList[0]
	return &shortening, nil
}

func (r *Repo) GetByURL(ctx context.Context, originURL string) (*dto.Shortening, error) {
	var shorteningList []dto.Shortening

	err := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.Select("short_url", "correlation_id", "original_url", "user_id").
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
			Columns("short_url", "correlation_id", "original_url", "user_id").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		for _, shortening := range list {
			sb = sb.Values(
				shortening.ShortURL,
				shortening.CorrelationID,
				shortening.OriginalURL,
				shortening.UserID,
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

func (r *Repo) GetListByUser(ctx context.Context, userID string) (*dto.ShorteningList, error) {
	var shorteningList dto.ShorteningList

	err := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.Select("short_url", "correlation_id", "original_url", "user_id").
			From("shorten_urls").
			PlaceholderFormat(squirrel.Dollar).
			Where(squirrel.Eq{"user_id": []string{userID}}).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &shorteningList, query, args...)
	})
	log := logger.LogFromContext(ctx)
	if err != nil {
		log.Info("repository.GetListByUser", zap.Error(err))
		return nil, errors.Wrap(err, "repository.GetListByUser")
	}

	if len(shorteningList) == 0 {
		log.Info("repository.GetListByUser", zap.Error(sql.ErrNoRows))
		return nil, errors.Wrap(sql.ErrNoRows, "repository.GetListByUser")
	}
	return &shorteningList, nil
}

func (r *Repo) DeleteUserURL(ctx context.Context, shortList []string, userID int) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := prepareList(doneCh, shortList)
	channels := r.split(ctx, doneCh, inputCh)
	resultCh := merge(doneCh, channels...)

	log := logger.LogFromContext(ctx)
	err := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		toSQL, _, err := squirrel.StatementBuilder.
			Update("shorten_urls").
			Set("is_deleted", true).
			Where(squirrel.And{
				squirrel.Eq{"user_id": ""},
				squirrel.Eq{"short_url": ""}}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		stmt, err := r.db.PrepareContext(ctx, toSQL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for Shortening := range resultCh {
			_, err = stmt.ExecContext(ctx, true, strconv.Itoa(userID), Shortening.ShortURL)
			if err != nil {
				log.Debug("repository.DeleteUserUrl Error", zap.Error(err))
			}
		}
		return nil
	})
	if err != nil {
		log.Debug("repository.DeleteUserUrl Error", zap.Error(err))
	}
}

func prepareList(doneCh chan struct{}, input []string) <-chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			select {
			case <-doneCh:
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}

func (r *Repo) searchURLs(ctx context.Context, doneCh chan struct{}, inputCh <-chan string) <-chan *dto.Shortening {
	addRes := make(chan *dto.Shortening)

	go func() {
		defer close(addRes)

		for shortURL := range inputCh {
			shortening, _ := r.Get(ctx, shortURL)

			select {
			case <-doneCh:
				return
			case addRes <- shortening:
			}
		}
	}()
	return addRes
}

func (r *Repo) split(ctx context.Context, doneCh chan struct{}, inputCh <-chan string) []<-chan *dto.Shortening {
	numWorkers := 100
	channels := make([]<-chan *dto.Shortening, numWorkers)

	for i := 0; i < numWorkers; i++ {
		addResultCh := r.searchURLs(ctx, doneCh, inputCh)
		channels[i] = addResultCh
	}

	return channels
}

func merge(doneCh chan struct{}, resultChs ...<-chan *dto.Shortening) <-chan *dto.Shortening {
	finalCh := make(chan *dto.Shortening)

	var wg sync.WaitGroup

	for _, ch := range resultChs {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-doneCh:
					return
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
