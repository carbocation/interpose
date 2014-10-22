/*
Interpose is a minimalist Golang middleware that uses only http.Handler and
func(http.Handler)http.Handler . Interpose takes advantage of closures to create
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

Therefore, the last middleware generator to be added will not only be
the first to be called, but will also have the opportunity to make the
final call after the rest of the middleware is called
*/
package interpose

import (
	"net/http"
)

type Middleware struct {
	Wares []func(http.Handler) http.Handler
}

// Return an empty middleware that is ready to use
func New() *Middleware {
	return &Middleware{}
}

// Add a piece of middleware which is an http.Handler generator
// (signature: func(http.Handler)http.Handler) which, somewhere before it
// finishes, is expected to call .ServeHTTP on the handler that is passed to it.
// Failure to call .ServeHTTP within the http.Handler generator will cause part
// of the stack not to be called.
func (mw *Middleware) Use(handler func(http.Handler) http.Handler) {
	mw.Wares = append(mw.Wares, handler)
}

// Add a piece of middleware which is simply any http.Handler
// (signature: http.Handler). Unlike with Use, we will automatically call
// .ServeHTTP to ensure that the rest of the middleware stack is called.
func (mw *Middleware) UseHandler(handler http.Handler) {
	x := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(w, req)
			next.ServeHTTP(w, req)
		})
	}

	mw.Use(x)
}

// Satisfies the net/http Handler interface and calls the middleware stack
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if len(mw.Wares) < 1 {
		return
	}

	//Initialize with an empty http.Handler
	next := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {}))

	//Call the middleware stack in FIFO order
	for i := len(mw.Wares) - 1; i >= 0; i-- {
		next = mw.Wares[i](next)
	}

	//Finally, serve back up the chain
	next.ServeHTTP(w, req)
}
