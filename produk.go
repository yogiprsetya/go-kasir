package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "Kecap", Harga: 23000, Stok: 41},
}

// Handle /products endpoint
func handleProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(produk)
	case "POST":
		var newProduct Produk
		err := json.NewDecoder(r.Body).Decode(&newProduct)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newProduct.ID = len(produk) + 1
		produk = append(produk, newProduct)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newProduct)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET api/produk/:id
// Handle /products/{id} endpoint
func handleProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Path[len("/api/produk/"):]
	id, err := strconv.Atoi(idStr)

	switch r.Method {
	case "GET":
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		// Find product
		for _, p := range produk {
			if p.ID == id {
				json.NewEncoder(w).Encode(p)
				return
			}
		}

		http.Error(w, "Product not found", http.StatusNotFound)
	case "PUT":
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var updatedProduct Produk
		err = json.NewDecoder(r.Body).Decode(&updatedProduct)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Find and update product
		for i, p := range produk {
			if p.ID == id {
				updatedProduct.ID = id
				produk[i] = updatedProduct
				json.NewEncoder(w).Encode(updatedProduct)
				return
			}
		}

		http.Error(w, "Product not found", http.StatusNotFound)
	case "DELETE":
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		// Find and delete product
		for i, p := range produk {
			if p.ID == id {
				produk = append(produk[:i], produk[i+1:]...)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Product deleted successfully",
				})
				return
			}
		}

		http.Error(w, "Product not found", http.StatusNotFound)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
