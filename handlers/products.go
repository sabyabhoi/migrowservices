package handlers

import (
	"log"
	"net/http"
	"regexp"

	"github.com/sabyabhoi/microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProducts(w, r)
		return
	}

	if r.Method == http.MethodPut {
		re := regexp.MustCompile(`/([0-9]+)`)
		group := re.FindAllStringSubmatch(r.URL.Path, -1)

		if len(group) != 1 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()

	w.Header().Add("Content-Type", "application/json")

	err := products.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to Marshall JSON", http.StatusInternalServerError)
	}
}

func (p *Products) addProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post requests")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to Unmarshall JSON", http.StatusBadRequest)
		return
	}
	data.AddProduct(prod)
}
