package main

import (
	"context"
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
	flagConfig := config.NewConfig()
	config.ParseFlags(flagConfig)

	ctx := config.ContextWithConfig(context.Background(), flagConfig)

	if err := app.Run(ctx); err != nil {
		fmt.Println(err)
	}
}
