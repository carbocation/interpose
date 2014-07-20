package middleware

import (
	"net/http"
)

/*
Middleware that sends an application/json header
*/
func Json(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}
