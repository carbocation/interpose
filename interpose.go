/*
Interpose is a minimalist Golang middleware that uses only http.Handler and
func(http.Handler)http.Handler taking advantage of closures to create a stack
of native net/http middleware that doesn't break, smudge, or otherwise require
any addition to the http.Handler interface.

From the view of a sandwich, the first middleware that you add gets called in
the middle, while the last middleware that you add gets called first (and can
make additional calls after earlier middleware finishes).
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

// Add a pice of middleware which is simply any http.Handler
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

func (mw *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if len(mw.Wares) < 1 {
		return
	}

	//Initialize with an empty http.Handler
	next := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {}))

	/*
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
	*/
	for _, generate := range mw.Wares {
		next = generate(next)
	}

	//Finally, serve back up the chain
	next.ServeHTTP(w, req)
}
