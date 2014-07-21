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
