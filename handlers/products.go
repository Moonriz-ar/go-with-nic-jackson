package handlers

import (
	"context"
	"fmt"
	"learn-go/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// getProducts returns the products from the data store
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle get products")
	listProducts := data.GetProducts()

	// fetch products from the datastore
	err := listProducts.ToJSON(w)
	// serialiaze the list to JSON
	if err != nil {
		http.Error(w, "unable to marshal json", http.StatusInternalServerError)
	}
}

// addProduct creates a new product
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle post product")

	product := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&product)
}

// updateProduct updates a product by id
func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
	}

	p.l.Println("handle update product")

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)

	if err == data.ErrProductNotFound {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(w, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				w,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		// call the next handler, which can be another middle in the chain, or the final handler
		next.ServeHTTP(w, req)
	})
}
