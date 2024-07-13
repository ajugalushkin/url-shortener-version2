//go:build windows || linux

package main

import (
	"context"
	"fmt"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/app"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
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

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	if err := app.Run(ctx); err != nil {
		fmt.Println(err)
	}
}
