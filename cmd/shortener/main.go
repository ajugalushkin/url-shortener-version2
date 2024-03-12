package main

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/app"
	"github.com/ajugalushkin/url-shortener-version2/internal/config"
)

func main() {
	config.ParseFlags()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
