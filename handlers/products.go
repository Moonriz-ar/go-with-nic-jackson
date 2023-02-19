package handlers

import (
	"learn-go/data"
	"log"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(w, r)
		return
	}

	// catch all
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	listProducts := data.GetProducts()

	// 1/ use json.Marshal to return a json of products.
	/*
		json, err := json.Marshal(listProducts)
		if err != nil {
			http.Error(w, "unable to marshal json", http.StatusInternalServerError)
		}
		w.Write(json)
	*/

	// 2/ use json.NewEncoder, does not allocate memory for marshal json
	// a little bit more efficient, especially if json is very big
	err := listProducts.ToJSON(w)
	if err != nil {
		http.Error(w, "unable to marshal json", http.StatusInternalServerError)
	}
}
