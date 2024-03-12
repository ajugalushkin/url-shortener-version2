package save

import (
	"fmt"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"io"
	"net/http"
)

func New(serviceAPI *service.Service) http.HandlerFunc {
	return func(wrt http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(wrt, "Invalid request method", http.StatusBadRequest)
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

		shortenURL, err := serviceAPI.Shorten(model.ShortenInput{RawURL: originalURL})
		if err != nil {
			http.Error(wrt, "URL not shortening", http.StatusBadRequest)
			return
		}

		shortenedURL := fmt.Sprintf("http://localhost:8080/%s", shortenURL.Key)

		wrt.Header().Set("Content-Type", "text/plain")
		wrt.WriteHeader(http.StatusCreated)
		wrt.Write([]byte(shortenedURL))
	}
}
