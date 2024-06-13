package config

import (
	"context"
	"flag"
	"os"
)

type Config struct {
	RunAddr         string `env:"RUN_ADDR"`
	BaseURL         string `env:"BASE_URL"`
	FlagLogLevel    string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDsn     string `env:"DATABASE_DSN"`
	SecretKey       string `env:"SECRET_KEY"`
}

func NewConfig() *Config {
	return &Config{}
}

func ParseFlags(config *Config) {
	flag.StringVar(&config.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base URL for POST request")
	flag.StringVar(&config.FlagLogLevel, "l", "info", "Log level")
	flag.StringVar(&config.FileStoragePath, "f", "",
		"full name of the file where data in JSON format is saved")
	flag.StringVar(&config.DataBaseDsn, "d", "",
		"DB path for connect")
	flag.Parse()

	config.RunAddr = getEnv("RUN_ADDR", config.RunAddr)
	config.BaseURL = getEnv("BASE_URL", config.BaseURL)
	config.FlagLogLevel = getEnv("LOG_LEVEL", config.FlagLogLevel)
	config.FileStoragePath = getEnv("FILE_STORAGE_PATH", config.FileStoragePath)
	config.DataBaseDsn = getEnv("DATABASE_DSN", config.DataBaseDsn)

	config.SecretKey = getEnv("SECRET_KEY", config.SecretKey)
	if config.SecretKey == "" {
		config.SecretKey = "SecretKey"
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

type ctxConfig struct{}

func ContextWithFlags(ctx context.Context, config *Config) context.Context {
	return context.WithValue(ctx, ctxConfig{}, config)
}

func FlagsFromContext(ctx context.Context) *Config {
	if config, ok := ctx.Value(ctxConfig{}).(*Config); ok {
		return config
	}
	return &Config{}
}
