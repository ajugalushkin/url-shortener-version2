package main

import (
	"fmt"
	"github.com/ajugalushkin/url-shortener-version2/internal/app"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
)

// @title shortener-url API
// @version 1.0
// @description Shorting URL API

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	cfg := config.NewConfig()
	config.ParseFlags(cfg)

	if err := app.Run(cfg); err != nil {
		fmt.Println(err)
	}
}
