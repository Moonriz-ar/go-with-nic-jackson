package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("hello world")

	// read request body
	d, err := io.ReadAll(r.Body)

	// if there is error while reading req body, server response with http error 400 and error message
	if err != nil {
		http.Error(w, "Oops", http.StatusBadRequest)
		// terminate the flow
		return
	}

	// write to response
	fmt.Fprintf(w, "Hello %s\n", d)
}
