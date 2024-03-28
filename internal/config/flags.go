package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr         string `env:"RUN_ADDR"`
	BaseURL         string `env:"BASE_URL"`
	FlagLogLevel    string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() *Config {
	return &Config{}
}

func ParseFlags(config *Config) {
	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		config.RunAddr = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.BaseURL = envBaseURL
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		config.FlagLogLevel = envLogLevel
	}

	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		config.FileStoragePath = envStoragePath
	}
	flag.StringVar(&config.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base URL for POST request")
	flag.StringVar(&config.FlagLogLevel, "l", "info", "Log level")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/short-url-db.json",
		"full name of the file where data in JSON format is saved")
	flag.Parse()
}
