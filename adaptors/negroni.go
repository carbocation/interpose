package adaptors

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

type NegroniRW struct {
	http.ResponseWriter
	//http.Flusher
}

func (nrw NegroniRW) Status() int   { return 0 }
func (nrw NegroniRW) Written() bool { return false }
func (nrw NegroniRW) Size() int     { return 0 }
func (nrw NegroniRW) Before(before func(negroni.ResponseWriter)) {
	//nrw.beforeFuncs = append(nrw.beforeFuncs, before)
}
func (nrw NegroniRW) Flush() {
	flusher, ok := nrw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

func FromNegroni(handler negroni.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		newRW := negroni.ResponseWriter(NegroniRW{rw})

		handler.ServeHTTP(newRW, req, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}))
	})
}
