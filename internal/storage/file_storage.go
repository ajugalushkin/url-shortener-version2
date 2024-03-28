package storage

import (
	"errors"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"sync"
)

type FileStorage struct {
	m   sync.Map
	cfg *config.Config
}

func NewStorage(cfg *config.Config) *FileStorage {
	storage := FileStorage{cfg: cfg}

	_ = refreshMap(&storage.m, storage.cfg.FileStoragePath)

	return &storage
}

func (s *FileStorage) Put(shortening model.Shortening) (*model.Shortening, error) {
	if _, exists := s.m.Load(shortening.Key); exists {
		return nil, errors.New("identifier already exists")
	}

	Writer, err := NewWriter(s.cfg.FileStoragePath)
	if err != nil {
		return nil, err
	}
	defer Writer.Close()

	if err := Writer.WriteFile(&model.File{
		ShortURL:    shortening.Key,
		OriginalURL: shortening.URL}); err != nil {
		return nil, err
	}

	if err := refreshMap(&s.m, s.cfg.FileStoragePath); err != nil {
		return nil, err
	}

	return &shortening, nil
}

func (s *FileStorage) Get(identifier string) (*model.Shortening, error) {
	v, ok := s.m.Load(identifier)
	if !ok {
		return nil, errors.New("not found")
	}

	shortening := v.(model.Shortening)

	return &shortening, nil
}

func refreshMap(m *sync.Map, filePath string) error {
	Reader, err := NewReader(filePath)
	if err != nil {
		return err
	}
	defer Reader.Close()

	files, err := Reader.ReadFile()
	if err != nil {
		return err
	}

	//erase SyncMap
	m.Range(func(key interface{}, value interface{}) bool {
		m.Delete(key)
		return true
	})

	for _, file := range files {
		m.Store(file.ShortURL,
			model.Shortening{
				Key: file.ShortURL,
				URL: file.OriginalURL})
	}

	return nil
}
