package routes

import (
	"net/http"
	"strings"

	"github.com/BerylCAtieno/string-analyzer/handlers"
)

// RegisterRoutes configures all API routes for the string analyzer
func RegisterRoutes() {
	// 1. Natural Language Filter Route (Most specific - register first)
	http.HandleFunc("/strings/filter-by-natural-language", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.NaturalLanguageFilterHandler(w, r)
	})

	// 2. Explicit handler for /strings (collection only - POST, GET)
	http.HandleFunc("/strings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// POST /strings: Create a new string
			handlers.CreateStringHandler(w, r)
		case http.MethodGet:
			// GET /strings: Get all strings (with or without filters)
			handlers.GetAllStringsHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 3. Handler for /strings/ (individual items - GET, DELETE)
	http.HandleFunc("/strings/", func(w http.ResponseWriter, r *http.Request) {
		// Extract the value after /strings/
		pathValue := strings.TrimPrefix(r.URL.Path, "/strings/")

		// If pathValue is empty, it means the request was to /strings/ with nothing after
		// This shouldn't reach here if /strings is properly registered, but handle it
		if pathValue == "" {
			http.Error(w, "string value not provided", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// GET /strings/{value}: Get a specific string
			handlers.GetStringHandler(w, r)

		case http.MethodDelete:
			// DELETE /strings/{value}: Delete a specific string
			handlers.DeleteStringHandler(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 4. Health check route
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}
