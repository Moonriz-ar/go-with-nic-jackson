package main

import (
	"context"
	"learn-go/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// create a new logger
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// create handlers
	productHandler := handlers.NewProducts(l)

	// create new ServeMux and register handlers
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productHandler.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", productHandler.UpdateProduct)
	putRouter.Use(productHandler.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productHandler.AddProduct)
	postRouter.Use(productHandler.MiddlewareProductValidation)

	// start and configure server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections.
	go func() {
		l.Println("starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	sig := <-signals
	l.Println("received terminate, graceful shutdown", sig)

	duration := time.Now().Add(30 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), duration)
	defer cancel()
	s.Shutdown(ctx)
}
