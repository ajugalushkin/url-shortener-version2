package main

import (
	"fmt"
	"github.com/ajugalushkin/url-shortener-version2/internal/app"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
)

func main() {
	cfg := config.NewConfig()
	config.ParseFlags(cfg)

	if err := app.Run(cfg); err != nil {
		fmt.Println(err)
	}
}
