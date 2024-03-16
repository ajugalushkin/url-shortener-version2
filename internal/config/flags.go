package config

import (
	"flag"
	"os"
)

type Config struct {
	FlagRunAddr string
	BaseURL     string
}

func (c Config) ParseFlags() {
	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		c.FlagRunAddr = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		c.BaseURL = envBaseURL
	}

	flag.StringVar(&c.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Base URL for POST request")
	flag.Parse()
}
