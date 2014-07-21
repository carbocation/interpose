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
	greenRouter := mux.NewRouter().Methods("GET").PathPrefix("/green").Subrouter() //Headers("Accept", "application/json")
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
