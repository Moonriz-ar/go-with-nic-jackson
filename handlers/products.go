package handlers

import (
	"learn-go/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle the request for a list a products
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	// handle the request to create a new product
	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("PUT")
		// expect the id in the URI
		regex := regexp.MustCompile(`/([0-9]+)`)
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Printf("invalid uri more than one id, got: %v", g)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Printf("invalid uri more than one capture group, got: %v", g)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Printf("invalid uri unable to convert to int, got: %v", idString)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, w, r)
		return
	}

	// catch all
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts returns the products from the data store
func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
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
func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle post product")

	product := &data.Product{}

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(product)
}

// updateProduct updates a product by id
func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle update product")

	product := &data.Product{}

	if err := product.FromJSON(r.Body); err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
	}

	err := data.UpdateProduct(id, product)

	if err == data.ErrProductNotFound {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "product not found", http.StatusInternalServerError)
		return
	}
}
