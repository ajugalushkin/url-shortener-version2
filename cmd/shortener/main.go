package main

import "github.com/ajugalushkin/url-shortener-version2/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
