interpose
=========

Interpose is a minimalist net/http middleware for golang. It uses 
http.Handler as its core unit of functionality, minimizing complexity
and maximizing inter-operability with other middleware frameworks.

Because of its reliance on the net/http standard, Interpose is out-of-the-box 
compatible with the Gorilla framework, goji, nosurf, and many other frameworks and 
standalone middleware.

A goal of the project is also to create adaptors so that non-http.Handler 
compliant middleware can still be used. As an example of this, an adaptor 
for Negroni middleware is available, making any middleware that is 
Negroni compliant also Interpose compliant. 

To use, first:

`go get github.com/carbocation/interpose`

## basic usage example

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

## Philosophy

Interpose is a minimalist Golang middleware that uses only http.Handler and
func(http.Handler)http.Handler . Interpose takes advantage of closures to create
a stack of native net/http middleware. Unlike other middleware libraries which
create their own net/http-like signatures, interpose uses literal net/http
signatures, thus minimizing package lock-in and maximizing inter-compatibility.

From the view of a sandwich, the first middleware that you add gets called in
the middle, while the last middleware that you add gets called first (and can
make additional calls after earlier middleware finishes).

The actual call stack of our chain of handlers starts from the last
added and ends with the first added. For example, if there are 3
middlewares added in order (0, 1, 2), the calls look like so:

	//2 START
		//1 START
			//0 START
			//0 END
		//1 END
	//2 END

Therefore, the last middleware generator to be added will not only be
the first to be called, but will also have the opportunity to make the
final call after the rest of the middleware is called

## Middleware

Here is a current list of Interpose compatible middleware. Feel free to put up a PR linking your middleware if you have built one:


| Middleware | Author | Description |
| -----------|--------|-------------|
| [Graceful](https://github.com/stretchr/graceful) | [Tyler Bunnell](https://github.com/tylerb) | Graceful HTTP Shutdown |
| [secure](https://github.com/unrolled/secure) | [Cory Jacobsen](https://github.com/unrolled) | Middleware that implements a few quick security wins |
| [logrus](https://github.com/carbocation/interpose/blob/master/examples/adaptors/logrus/main.go) | [Dan Buch](https://github.com/meatballhat) | Logrus-based logger demonstrating how Negroni packages can be used in Interpose |
| [buffer](https://github.com/carbocation/interpose/blob/master/examples/buffer/main.go) | [carbocation](https://github.com/carbocation) | Output buffering demonstrating how headers can be written after HTTP body is sent |