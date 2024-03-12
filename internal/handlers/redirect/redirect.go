package redirect

import (
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"net/http"
	"strings"
)

func New(serviceAPI *service.Service) http.HandlerFunc {
	return func(wrt http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(wrt, "Invalid request method", http.StatusBadRequest)
			return
		}

		key := strings.Replace(req.URL.Path, "/", "", -1)

		redirect, err := serviceAPI.Redirect(key)
		if err != nil {
			http.Error(wrt, "Original URL not found!", http.StatusBadRequest)
			return
		}

		wrt.Header().Set("Location", redirect)
		wrt.WriteHeader(http.StatusTemporaryRedirect)
	}
}
