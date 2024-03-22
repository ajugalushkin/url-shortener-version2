package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr      string `env:"RUN_ADDR"`
	BaseURL      string `env:"BASE_URL"`
	FlagLogLevel string `env:"LOG_LEVEL"`
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

	flag.StringVar(&config.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base URL for POST request")
	flag.StringVar(&config.FlagLogLevel, "c", "info", "Log level")
	flag.Parse()
}
