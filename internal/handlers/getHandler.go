package handlers

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/storage"
	"net/http"
	"strings"
)

func GetHandler(wrt http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(wrt, "Invalid request method", http.StatusBadRequest)
		return
	}

	key := strings.Replace(req.URL.Path, "/", "", -1)

	storageAPI, errGetAPI := storage.NewStorage()
	if errGetAPI != nil {
		http.Error(wrt, "Storage not found!", http.StatusBadRequest)
		return
	}

	dataURL, err := storageAPI.Retrieve(key)
	if err != nil {
		http.Error(wrt, err.Error(), http.StatusBadRequest)
		return
	}

	wrt.Header().Set("Location", dataURL.URL)
	wrt.WriteHeader(http.StatusTemporaryRedirect)
}
