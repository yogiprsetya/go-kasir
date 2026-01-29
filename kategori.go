package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Sesuatu yang bisa dimakan"},
	{ID: 2, Name: "Minuman", Description: "Sesuatu yang bisa diminum"},
}

// Handle /categories endpoint
func handleCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(categories)
	case "POST":
		var newCategory Category
		err := json.NewDecoder(r.Body).Decode(&newCategory)

		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newCategory.ID = len(categories) + 1
		categories = append(categories, newCategory)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCategory)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// / Handle /categories/{id} endpoint
func handleCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		idStr := r.URL.Path[len("/api/kategori/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		for _, c := range categories {
			if c.ID == id {
				json.NewEncoder(w).Encode(c)
				return
			}
		}

		http.Error(w, "Category not found", http.StatusNotFound)
	case "PUT":
		idStr := r.URL.Path[len("/api/kategori/"):]
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var updatedCategory Category
		err = json.NewDecoder(r.Body).Decode(&updatedCategory)

		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		for i, c := range categories {
			if c.ID == id {
				updatedCategory.ID = id
				categories[i] = updatedCategory
				json.NewEncoder(w).Encode(updatedCategory)
				return
			}
		}

		http.Error(w, "Category not found", http.StatusNotFound)
	case "DELETE":
		idStr := r.URL.Path[len("/api/kategori/"):]
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		for i, c := range categories {
			if c.ID == id {
				categories = append(categories[:i], categories[i+1:]...)

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Category deleted successfully",
				})
				return
			}
		}

		http.Error(w, "Category not found", http.StatusNotFound)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
