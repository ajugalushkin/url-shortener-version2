package repository

import (
	"context"
	"database/sql"
	"strconv"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ajugalushkin/url-shortener-version2/internal/database"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	userErr "github.com/ajugalushkin/url-shortener-version2/internal/errors"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
)

// NewRepository Конструктор
func NewRepository(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// Repo структура репозитория
type Repo struct {
	db *sqlx.DB
}

// Put метод сохраняет данные URL в базу данных.
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

// Get метод позволяет получить оригинальный URL по сокращенному.
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

	log := logger.GetLogger()
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

// GetByURL метод позволяет получить данные URL по оригинальному URL.
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

	log := logger.GetLogger()
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

// PutList метод позволяет сохранить список URL в базу данных.
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

// GetListByUser метод получает все сокращенные URL для конкретного пользователя.
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

	log := logger.GetLogger()
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

// DeleteUserURL метод удаляет все сокращенные URL для конкретного пользователя.
func (r *Repo) DeleteUserURL(ctx context.Context, shortList []string, userID int) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := prepareList(doneCh, shortList)
	channels := r.split(ctx, doneCh, inputCh)
	resultCh := merge(doneCh, channels...)

	log := logger.GetLogger()
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

// prepareList функция реализует Fan-Out позволяет передавать входящий список URL в список каналов.
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

// searchURLs метод используется для поиска данных URL, используется в рамках Fan-Out
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

// split метод реализует Fan-Out позволяет назначмть воркеров для обработки входящих данных через канал.
func (r *Repo) split(ctx context.Context, doneCh chan struct{}, inputCh <-chan string) []<-chan *dto.Shortening {
	numWorkers := 100
	channels := make([]<-chan *dto.Shortening, numWorkers)

	for i := 0; i < numWorkers; i++ {
		addResultCh := r.searchURLs(ctx, doneCh, inputCh)
		channels[i] = addResultCh
	}

	return channels
}

// merge функция реализует Fan-In позволяет получить итоговый результат из канала.
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

// GetStats метод для получения статистики
func (r *Repo) GetStats(ctx context.Context) (*dto.Stats, error) {
	var stats dto.Stats
	err := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.Select("count(short_url) AS urls", "count(user_id) as users").
			From("shorten_urls").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &stats, query, args...)
	})

	if err != nil {
		logger.GetLogger().Info("repository.GetStats", zap.Error(err))
		return nil, errors.Wrap(err, "repository.GetStats")
	}

	log.Info(
		"repository.GetStats OK",
		zap.String("URLS", strconv.Itoa(stats.URLS)),
		zap.String("Users", strconv.Itoa(stats.Users)),
	)

	return &stats, nil
}
