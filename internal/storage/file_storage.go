package storage

import (
	"bytes"
	"errors"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage struct {
	m   sync.Map
	cfg *config.Config
}

func NewStorage(cfg *config.Config) *FileStorage {
	storage := FileStorage{cfg: cfg}
	_ = load(&storage.m, storage.cfg.FileStoragePath)
	return &storage
}

func (s *FileStorage) Put(shortening model.Shortening) (*model.Shortening, error) {
	if _, exists := s.m.Load(shortening.Key); exists {
		return nil, errors.New("identifier already exists")
	}

	s.m.Store(shortening.Key, shortening)

	err := save(s.cfg.FileStoragePath, &s.m)
	if err != nil {
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

func save(fileName string, urls *sync.Map) error {
	var byteFile []byte
	urls.Range(func(k, v interface{}) bool {
		shortening := v.(model.Shortening)

		file := model.File{
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
		file := model.File{}
		err := file.UnmarshalJSON(item)
		if err != nil {
			return err
		}
		files.Store(file.ShortURL, model.Shortening{Key: file.ShortURL, URL: file.OriginalURL})
	}

	return nil
}
