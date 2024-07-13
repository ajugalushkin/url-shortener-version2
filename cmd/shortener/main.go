package main

import (
	"context"
	"fmt"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/app"
)

// @title shortener-url API
// @version 1.0
// @description Shorting URL API

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	cfg := config.ReadConfig()
	ctx := config.ContextWithFlags(context.Background(), cfg)

	if err := app.Run(ctx); err != nil {
		fmt.Println(err)
	}
}
