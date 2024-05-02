package file

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
)

type Storage struct {
	m sync.Map
}

func NewStorage(path string) *Storage {
	storage := Storage{}
	_ = load(&storage.m, path)
	return &storage
}

func (r *Storage) Put(ctx context.Context, shortening dto.Shortening) (*dto.Shortening, error) {
	if _, exists := r.m.Load(shortening.ShortURL); exists {
		return nil, errors.New("identifier already exists")
	}

	r.m.Store(shortening.ShortURL, shortening)

	flags := config.FlagsFromContext(ctx)
	err := save(flags.FileStoragePath, &r.m)
	if err != nil {
		return nil, err
	}

	return &shortening, nil
}

func (r *Storage) PutList(ctx context.Context, list dto.ShorteningList) error {
	for _, shortening := range list {
		_, err := r.Put(ctx, shortening)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Storage) Get(ctx context.Context, identifier string) (*dto.Shortening, error) {
	v, ok := r.m.Load(identifier)
	if !ok {
		return nil, errors.New("not found")
	}

	shortening := v.(dto.Shortening)

	return &shortening, nil
}

func (r *Storage) GetListByUser(ctx context.Context, userID string) (*dto.ShorteningList, error) {
	return &dto.ShorteningList{}, nil
}

func (r *Storage) DeleteUserURL(ctx context.Context, shortURL string, userID int) error {
	return nil
}

func save(fileName string, urls *sync.Map) error {
	var byteFile []byte
	urls.Range(func(k, v interface{}) bool {
		shortening := v.(dto.Shortening)

		file := dto.Shortening{
			ShortURL:    shortening.ShortURL,
			OriginalURL: shortening.OriginalURL}

		data, err := file.MarshalJSON()
		if err != nil {
			return false
		}
		data = append(data, '\n')
		byteFile = append(byteFile, data...)

		return true
	})

	fileName = filepath.FromSlash(fileName)
	directory, _ := filepath.Split(fileName)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(fileName, byteFile, 0666)
}

func load(files *sync.Map, fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	splitData := bytes.Split(data, []byte("\n"))

	for _, item := range splitData {
		file := dto.Shortening{}
		err := file.UnmarshalJSON(item)
		if err != nil {
			return err
		}
		files.Store(file.ShortURL, dto.Shortening{ShortURL: file.ShortURL, OriginalURL: file.OriginalURL})
	}

	return nil
}
