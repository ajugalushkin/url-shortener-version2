package handlers

import (
	"fmt"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"github.com/ajugalushkin/url-shortener-version2/internal/shorten"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"io"
	"net/http"
)

func PostHandler(wrt http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		GetHandler(wrt, req)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(wrt, "URL parse error", http.StatusBadGateway)
		return
	}

	originalURL := string(body)
	if originalURL == "" {
		http.Error(wrt, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	storageAPI, errGetAPI := storage.NewStorage()
	if errGetAPI != nil {
		http.Error(wrt, "Storage not found!", http.StatusBadRequest)
		return
	}

	URLData, errGet := storageAPI.RetrieveByURL(originalURL)
	if errGet != nil {
		URLData = model.URLData{
			Key: shorten.GenerateShortKey(),
			URL: originalURL}
		_, err = storageAPI.Insert(URLData)
		if err != nil {
			http.Error(wrt, "ShortKey not created", http.StatusNotFound)
			return
		}
	}

	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", URLData.Key)

	wrt.Header().Set("Content-Type", "text/plan")
	wrt.WriteHeader(http.StatusCreated)
	wrt.Write([]byte(shortenedURL))
}
