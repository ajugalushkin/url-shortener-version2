package config

import "flag"

var (
	FlagRunAddr string
	BaseURL     string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Base URL for POST request")
	flag.Parse()
}
