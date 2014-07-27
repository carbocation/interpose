package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
)

type key int

const (
	CountKey key = iota
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		c, ok := context.GetOk(req, CountKey)
		if !ok {
			fmt.Println("Get not ok")
		}

		fmt.Fprintf(w, "Welcome to the home page, %s!\nCount:%d", mux.Vars(req)["user"], c)

	})

	mw := interpose.New()

	// Apply the router.
	mw.UseHandler(router)

	mw.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			c := rand.Int()
			fmt.Println("Setting ctx count to:", c)
			context.Set(req, CountKey, c)
			next.ServeHTTP(w, req)
		})
	})
	mw.Use(context.ClearHandler)

	// Launch and permit graceful shutdown, allowing up to 10 seconds for existing
	// connections to end
	graceful.Run(":3001", 10*time.Second, mw)
}
