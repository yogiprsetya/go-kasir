package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Product routes
	http.HandleFunc("/api/produk", handleProducts)
	http.HandleFunc("/api/produk/", handleProductByID)

	// Category routes
	http.HandleFunc("/api/kategori", handleCategories)
	http.HandleFunc("/api/kategori/", handleCategoryByID)

	fmt.Println("Server running on :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
