package middleware

import (
	"net/http"

	"github.com/carbocation/interpose/adaptors"
	"github.com/phyber/negroni-gzip/gzip"
)

/*
func Gzip(compression int) http.Handler {
	return adaptors.HandlerFromNegroni(gzip.Gzip(compression))
}
*/

func Gzip(compression int) func(http.Handler) http.Handler {
	return adaptors.FromNegroni(gzip.Gzip(compression))
}

/*
func Gzip(compression int) func(http.Handler) http.Handler {
	this := adaptors.HandlerFromNegroni(gzip.Gzip(compression))
	return func(next http.Handler) http.Handler {
		return http.Handler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			this.ServeHTTP(rw, req)
			next.ServeHTTP(rw, req)
		}))
	}
}

*/
