package routes

import (
	"net/http"

	"github.com/BerylCAtieno/string-analyzer/handlers"
)

// RegisterRoutes configures all API routes for the string analyzer
func RegisterRoutes() {
	http.HandleFunc("/strings/filter-by-natural-language", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.NaturalLanguageFilterHandler(w, r)
	})

	http.HandleFunc("/strings/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetStringHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteStringHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/strings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateStringHandler(w, r)
		case http.MethodGet:
			handlers.GetAllStringsHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check route
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}
