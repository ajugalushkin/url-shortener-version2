package storage

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
	m   sync.Map
	ctx context.Context
}

func NewStorage(ctx context.Context) *Storage {
	storage := Storage{ctx: ctx}

	flags := config.ConfigFromContext(ctx)
	_ = load(&storage.m, flags.FileStoragePath)
	return &storage
}

func (s *Storage) Put(shortening dto.Shortening) (*dto.Shortening, error) {
	if _, exists := s.m.Load(shortening.Key); exists {
		return nil, errors.New("identifier already exists")
	}

	s.m.Store(shortening.Key, shortening)

	flags := config.ConfigFromContext(s.ctx)
	err := save(flags.FileStoragePath, &s.m)
	if err != nil {
		return nil, err
	}

	return &shortening, nil
}

func (s *Storage) Get(identifier string) (*dto.Shortening, error) {
	v, ok := s.m.Load(identifier)
	if !ok {
		return nil, errors.New("not found")
	}

	shortening := v.(dto.Shortening)

	return &shortening, nil
}

func save(fileName string, urls *sync.Map) error {
	var byteFile []byte
	urls.Range(func(k, v interface{}) bool {
		shortening := v.(dto.Shortening)

		file := dto.File{
			ShortURL:    shortening.Key,
			OriginalURL: shortening.URL}

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
		file := dto.File{}
		err := file.UnmarshalJSON(item)
		if err != nil {
			return err
		}
		files.Store(file.ShortURL, dto.Shortening{Key: file.ShortURL, URL: file.OriginalURL})
	}

	return nil
}
