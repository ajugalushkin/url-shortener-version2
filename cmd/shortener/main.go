package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"

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
	cfg := config.ReadConfig()
	ctx := config.ContextWithFlags(context.Background(), cfg)

	if err := app.Run(ctx); err != nil {
		fmt.Println(err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
}
