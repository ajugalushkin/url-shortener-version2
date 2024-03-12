package config

import (
	"flag"
	"os"
)

var (
	FlagRunAddr string
	BaseURL     string
)

func ParseFlags() {
	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		BaseURL = envBaseURL
	}

	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Base URL for POST request")
	flag.Parse()
}
