package middleware

import (
	"net/http"
)

func Switch(post http.Handler, get http.Handler) http.Handler {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			post.ServeHTTP(wrt, req)
		}
		if req.Method == http.MethodGet {
			get.ServeHTTP(wrt, req)
		}
	})
}
