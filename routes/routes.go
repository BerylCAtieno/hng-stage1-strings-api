package routes

import (
	"net/http"
	"strings"

	"github.com/BerylCAtieno/string-analyzer/handlers"
)

// RegisterRoutes configures all API routes for the string analyzer
func RegisterRoutes() {
	// 1. Natural Language Filter Route (Specific Path, remains isolated)
	http.HandleFunc("/strings/filter-by-natural-language", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.NaturalLanguageFilterHandler(w, r)
	})

	// 2. Consolidated Handler for /strings and /strings/{value}
	// Handles all GET, POST, DELETE operations for the resource.
	http.HandleFunc("/strings/", func(w http.ResponseWriter, r *http.Request) {
		// Clean the path to check if it's a collection operation or an item operation.
		// A request to /strings or /strings/ will result in an empty "value" after stripping the prefix.
		pathValue := strings.TrimPrefix(r.URL.Path, "/strings")
		isCollection := pathValue == "" || pathValue == "/"

		switch r.Method {
		case http.MethodPost:
			if !isCollection {
				// POST /strings/value is not supported
				http.Error(w, "method not allowed on item", http.StatusMethodNotAllowed)
				return
			}
			// POST /strings: Create a new string
			handlers.CreateStringHandler(w, r)

		case http.MethodGet:
			if isCollection {
				// GET /strings: Get all strings (with or without filters)
				handlers.GetAllStringsHandler(w, r)
			} else {
				// GET /strings/value: Get a specific string
				handlers.GetStringHandler(w, r)
			}

		case http.MethodDelete:
			if isCollection {
				// DELETE /strings: Mass delete is not supported
				http.Error(w, "method not allowed on collection", http.StatusMethodNotAllowed)
				return
			}
			// DELETE /strings/value: Delete a specific string
			handlers.DeleteStringHandler(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}
