package app

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/handlers/middleware"
	"github.com/ajugalushkin/url-shortener-version2/internal/handlers/redirect"
	"github.com/ajugalushkin/url-shortener-version2/internal/handlers/save"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"net/http"
)

func Run() error {
	mux := http.NewServeMux()

	serviceAPI := service.NewService(storage.NewInMemory())

	mux.Handle("/", middleware.Switch(save.New(serviceAPI), redirect.New(serviceAPI)))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		return err
	}
	return nil
}
