package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// registers a http handler function to root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello world")

		// read request body
		d, err := io.ReadAll(r.Body)

		// if there is error while reading req body, server response with http error 400 and error message
		if err != nil {
			http.Error(w, "Oops", http.StatusBadRequest)
			// terminate the flow
			return
		}

		// write to response
		fmt.Fprintf(w, "Hello %s", d)
	})

	// registers a http handler function a /goodbye path
	http.HandleFunc("/goodbye", func(w http.ResponseWriter, r *http.Request) {
		log.Println("goodbye world")
	})

	http.ListenAndServe(":9090", nil)
}
