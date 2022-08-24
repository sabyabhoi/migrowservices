// Package Classification Product API
//
// Documentation of Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sabyabhoi/microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// Returns a list of coffee products
// responses:
// 200: productsResponse

func (p Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()

	w.Header().Add("Content-Type", "application/json")

	err := products.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to Marshall JSON", http.StatusInternalServerError)
	}
}

func (p Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST requests")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

func (p Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id to int", http.StatusInternalServerError)
		return
	}

	p.l.Println("Handle PUT requests")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err != nil {
		http.Error(w, "Unable to update product", http.StatusBadRequest)
		return
	}
}

type KeyProduct struct {}

func (p Products) ValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Unable to Unmarshall JSON", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to Validate product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
