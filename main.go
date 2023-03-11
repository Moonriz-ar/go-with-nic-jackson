package main

import (
	"context"
	"learn-go/data"
	"learn-go/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// create a new logger
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	v := data.NewValidation()

	// create handlers
	productHandler := handlers.NewProducts(l, v)

	// create new ServeMux and register handlers
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productHandler.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", productHandler.Update)
	putRouter.Use(productHandler.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productHandler.Create)
	postRouter.Use(productHandler.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// cors
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:300`	0"}))

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

	// Block until a signal is received
	sig := <-signals
	l.Println("received terminate, graceful shutdown", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	duration := time.Now().Add(30 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), duration)
	defer cancel()
	s.Shutdown(ctx)
}
