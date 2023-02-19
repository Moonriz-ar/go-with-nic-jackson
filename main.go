package main

import (
	"context"
	"learn-go/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// create a new logger
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// inject logger dependency to handlers
	helloHandler := handlers.NewHello(l)
	goodbyeHandler := handlers.NewGoodbye(l)

	// create new ServeMux
	sm := http.NewServeMux()

	// register handler to path
	sm.Handle("/", helloHandler)
	sm.Handle("/goodbye", goodbyeHandler)

	// configure server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections.
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
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
