package middleware

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
)

// Basic returns a Handler that authenticates via Basic Auth. Writes a http.StatusUnauthorized
// if authentication fails
func BasicAuth(username string, password string) func(http.Handler) http.Handler {
	var siteAuth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			auth := req.Header.Get("Authorization")
			if !secureCompare(auth, "Basic "+siteAuth) {
				res.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
				http.Error(res, "Not Authorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(res, req)
		})
	}
}

// secureCompare performs a constant time compare of two strings to limit timing attacks.
func secureCompare(given string, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	} else {
		/* Securely compare actual to itself to keep constant time, but always return false */
		return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
	}
}
