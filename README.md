interpose
=========

Interpose is a minimalist net/http middleware for golang

To use, first:

`go get github.com/carbocation/interpose`

Here is one example of using Interpose along with gorilla/mux to create
middleware that adds JSON headers to every response.

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	mw := interpose.New()

	// Apply the router. By adding it first, all of our other middleware will be
	// executed before the router, allowing us to modify headers before any
	// output has been generated.
	mw.UseHandler(router)

	// Tell the browser our output will be JSON
	mw.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	})

	// Launch and permit graceful shutdown, allowing up to 10 seconds for existing
	// connections to end
	graceful.Run(":3001", 10*time.Second, mw)
}
```

## Middleware

Here is a current list of Interpose compatible middleware. Feel free to put up a PR linking your middleware if you have built one:


| Middleware | Author | Description |
| -----------|--------|-------------|
| [Graceful](https://github.com/stretchr/graceful) | [Tyler Bunnell](https://github.com/tylerb) | Graceful HTTP Shutdown |
| [secure](https://github.com/unrolled/secure) | [Cory Jacobsen](https://github.com/unrolled) | Middleware that implements a few quick security wins |
| [logrus](https://github.com/carbocation/interpose/examples/adaptors/logrus/logrus.go) | [Dan Buch](https://github.com/meatballhat) | Logrus-based logger demonstrating how Negroni packages can be used in Interpose |
| [buffer](https://github.com/carbocation/interpose/middleware/buffer/buffer.go) | [carbocation](https://github.com/carbocation) | Output buffering demonstrating how headers can be written after HTTP body is sent |
| [gzip](https://github.com/phyber/negroni-gzip) | [phyber](https://github.com/phyber) | GZIP response compression |