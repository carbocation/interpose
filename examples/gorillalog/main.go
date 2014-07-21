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

	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	// First apply any middleware that modify the http body, since the first
	// added will be the last applied. This permits other middleware to alter headers
	middle.UseHandler(router)

	// Now apply any middleware that will not write output to http body
	middle.Use(middleware.GorillaLog())

	http.ListenAndServe(":3001", middle)
}
