interpose
=========

Interpose is a minimalist net/http middleware framework for golang. It uses 
`http.Handler` as its core unit of functionality, minimizing complexity
and maximizing inter-operability with other middleware frameworks.

All that it does is manage middleware. It comes with nothing baked in. You 
bring your own router, etc. See below for some well-baked examples.

Because of its reliance on the net/http standard, Interpose is out-of-the-box 
compatible with the Gorilla framework, goji, nosurf, and many other frameworks and 
standalone middleware. 

Many projects claim to be `http.Handler`-compliant but actually just use `http.Handlers` 
to create a more complicated/less compatible abstraction. Therefore, a goal of the 
project is also to create adaptors so that non-`http.Handler` compliant middleware can 
still be used. As an example of this, an adaptor for Negroni middleware is available, 
meaning that **any middleware that is Negroni compliant is also Interpose compliant**. 
The same is true for Martini middleware.

To use, first:

`go get github.com/carbocation/interpose`

## basic usage example

Here is one example of using Interpose to execute 
middleware that adds the HTTP header "X-Server-Name" to every response.

Create a file `main.go` with the following:
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
)

func main() {
	// In this example, we use a basic router with a catchall route that
	// matches any path. In other examples in this project, we use a
	// more advanced router.
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to %s!", req.URL.Path)
	}))

	middle := interpose.New()

	// Apply the router. The first added is the last executed, so by adding
	// the router first, our other middleware can modify HTTP headers before our
	// router writes to the HTTP Body
	middle.UseHandler(router)

	// Send a header telling the world this is coming from an Interpose Test Server
	// You could imagine setting Content-type application/json or other useful
	// headers in other circumstances.
	middle.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("X-Server-Name", "Interpose Test Server")
	}))

	http.ListenAndServe(":3001", middle)
}

```

In the same path as that file, type `go run main.go`

Now launch your browser and point it to `http://localhost:3001/world` to see output.

Additional examples can be found below.

## Philosophy

Interpose is a minimalist Golang middleware that uses only `http.Handler` and
`func(http.Handler)http.Handler`. Interpose takes advantage of closures to create
a stack of native net/http middleware. Unlike other middleware libraries which
create their own net/http-like signatures, interpose uses literal net/http
signatures, thus minimizing package lock-in and maximizing inter-compatibility.

Middleware is called in nested LIFO fashion, which means that the last middleware
to be added will be the first middleware to be called. Because the middleware is
nested, it actually means that the last middleware to be added gets the
opportunity to make the first and the last calls in the stack. For example,
if there are 3 middlewares added in order (0, 1, 2), the calls look like so:

	//2 START
		//1 START
			//0 START
			//0 END
		//1 END
	//2 END

## Middleware

Here is a current list of Interpose compatible middleware that have pre-built 
examples working with Interpose. Any middleware that yields an `http.Handler` 
or a `func(http.Handler)http.Handler` should be compliant. Pull requests linking 
to other middleware are encouraged.


| Middleware | Usage example | Author | Description |
| -----------|---------------|--------|-------------|
| [Graceful](https://github.com/stretchr/graceful) | [Graceful example](https://github.com/carbocation/interpose/blob/master/examples/graceful/main.go) | [Tyler Bunnell](https://github.com/tylerb) | Graceful HTTP Shutdown |
| [secure](https://github.com/unrolled/secure) | [Secure example](https://github.com/carbocation/interpose/blob/master/examples/secure/main.go) | [Cory Jacobsen](https://github.com/unrolled) | Middleware that implements a few quick security wins |
| [Gorilla logger](https://github.com/gorilla/handlers) | [Gorilla log example](https://github.com/carbocation/interpose/blob/master/examples/gorillalog/main.go) | [Gorilla team](https://github.com/gorilla/) | Gorilla Apache CombinedLogger |
| [Logrus](https://github.com/meatballhat/negroni-logrus) | [Logrus example](https://github.com/carbocation/interpose/blob/master/examples/adaptors/logrus/main.go) | [Dan Buch](https://github.com/meatballhat) | Logrus-based logger, also demonstrating how Negroni packages can be used in Interpose |
| [Buffered output](https://github.com/goods/httpbuf) | [Buffer example](https://github.com/carbocation/interpose/blob/master/examples/buffer/main.go) | [zeebo](https://github.com/zeebo) | Output buffering demonstrating how headers can be written after HTTP body is sent |
| [nosurf](https://github.com/justinas/nosurf) | [nosurf example](https://github.com/carbocation/interpose/blob/master/examples/nosurf/main.go) | [justinas](https://github.com/justinas) | A CSRF protection middleware for Go. |
| [BasicAuth](https://github.com/carbocation/interpose/blob/master/middleware/basicAuth.go)| [BasicAuth example](https://github.com/carbocation/interpose/blob/master/examples/basicAuth/main.go)| [Jeremy Saenz](http://github.com/codegangsta) & [Brendon Murphy](http://github.com/bemurphy) | [HTTP BasicAuth](https://en.wikipedia.org/wiki/Basic_access_authentication) - based on martini's [auth](https://github.com/martini-contrib/auth) middleware|
| [Martini Auth](https://github.com/martini-contrib/auth) | [Martini Auth example](https://github.com/carbocation/interpose/blob/master/examples/adaptors/martiniauth/main.go) | [Jeremy Saenz](https://github.com/codegangsta) & [Brendon Murphy](http://github.com/bemurphy) | A basic HTTP Auth implementation that also demonstrates how Martini middleware packages can be used directly in Interpose with a simple wrapper. |

## Adaptors

Some frameworks that are not strictly `http.Handler` compliant use middleware that 
can be readily converted into Interpose-compliant middleware. So far, adaptors for 
Martini and Negroni have been created. 

For example, to use github.com/codegangsta/negroni middleware in Interpose, you 
can use `adaptors.FromNegroni`:

```go
	middle := interpose.New()

	// has signature `negroni.Handler`
	negroniMiddleware := negronilogrus.NewMiddleware()

	// Use the Negroni middleware within Interpose
	middle.Use(adaptors.FromNegroni(negroniMiddleware))

```

## More examples

### Routing, graceful shutdown, and JSON headers

In this example, we use the `graceful` package to gracefully release 
connections after the shutdown signal is encountered. A more powerful 
router than in our prior example, Gorilla mux, is used. Also, we send 
the browser headers indicating that this is a JSON response (which, 
however, it is not).

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

	middle := interpose.New()

	// Apply the router. By adding it first, all of our other middleware will be
	// executed before the router, allowing us to modify headers before any
	// output has been generated.
	middle.UseHandler(router)

	// Tell the browser which server this came from
	middle.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("X-Server-Name", "Interpose Test Server")
			next.ServeHTTP(rw, req)
		})
	})

	// Launch and permit graceful shutdown, allowing up to 10 seconds for existing
	// connections to end
	graceful.Run(":3001", 10*time.Second, middle)
}

```

### Combined logging and gzipping

Print an Apache CombinedLog-compatible log statement to StdOut and 
gzip the HTTP response it if the client has gzip capabilities:

```go
package main

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/mux"
)

func main() {
	middle := interpose.New()

	router := mux.NewRouter()
	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	// First apply any middleware that modify the http body, since the first
	// added will be the last applied. This permits other middleware to alter headers
	middle.UseHandler(router)

	// Now apply any middleware that will not write output to http body

	// Log to stdout. Taken from Gorilla
	middle.Use(middleware.GorillaLog())

	// Gzip output. Taken from Negroni
	middle.Use(middleware.NegroniGzip(gzip.DefaultCompression))

	http.ListenAndServe(":3001", middle)
}

```

### Wrapped middleware

Middleware can be wrapped around other middleware. In this example, 
we greet people who arrive at `/{user}`, but we offer a special greeting
to those who arrive at `/green/{user}`:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/mux"
)

func main() {
	middle := interpose.New()

	router := mux.NewRouter()
	router.Handle("/{name}", http.HandlerFunc(welcomeHandler))
	router.PathPrefix("/green").Subrouter().Handle("/{name}", Green(http.HandlerFunc(welcomeHandler)))

	middle.UseHandler(router)

	// Using Gorilla framework's combined logger
	middle.Use(middleware.GorillaLog())

	http.ListenAndServe(":3001", middle)
}

func welcomeHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Welcome to the home page, %s", mux.Vars(req)["name"])
}

func Green(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("X-Favorite-Color", "green")
		next.ServeHTTP(rw, req)
		fmt.Fprint(rw, " who likes green")
	})
}

```

### Nested middleware: adding headers for only some routes

In the last example, we applied different middleware to different routes. 
Here we will expand this idea to created fully nested middleware stacks 
within different routes. This approach, while more verbose than the last 
example, is arbitrarily powerful. In this example, routes starting with 
/green are given a special HTTP header X-Favorite-Color: green, 
but you can also imagine using this same approach to automatically apply 
the JSON content header for JSON requests, putting authentication in front of
protected paths, etc.

```go
package main

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/mux"
)

func main() {
	middle := interpose.New()

	router := mux.NewRouter()
	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	middle.UseHandler(router)

	// Using Gorilla framework's combined logger
	middle.Use(middleware.GorillaLog())

	//Using Negroni's Gzip functionality
	middle.Use(middleware.NegroniGzip(gzip.DefaultCompression))

	// Now we will define a sub-router based on our love of the color green
	// When you call any url such as http://localhost:3001/green/man , you will
	// also see an HTTP header sent called X-Favorite-Color with value "green"
	greenRouter := mux.NewRouter().Methods("GET").PathPrefix("/green").Subrouter()
	greenRouter.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, green %s!", mux.Vars(req)["user"])
	})

	greenMiddle := interpose.New()
	greenMiddle.UseHandler(greenRouter)
	greenMiddle.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("X-Favorite-Color", "green")
	}))
	router.Methods("GET").PathPrefix("/green").Handler(greenMiddle)

	http.ListenAndServe(":3001", middle)
}

```

For more examples, please look at the [examples folder](https://github.com/carbocation/interpose/tree/master/examples) 
as well as its subfolder, the [menagerie folder](https://github.com/carbocation/interpose/tree/master/examples/menagerie)

## Authors
Originally developed by [carbocation](https://github.com/carbocation). Please see the 
[contributors](https://github.com/carbocation/interpose/blob/master/CONTRIBUTORS.md) 
file for an expanded list of contributors.
