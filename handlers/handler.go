package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/BerylCAtieno/string-analyzer/analysis"
	"github.com/BerylCAtieno/string-analyzer/database"
	"github.com/BerylCAtieno/string-analyzer/filtering"
	"github.com/BerylCAtieno/string-analyzer/models"
	"github.com/BerylCAtieno/string-analyzer/nlp"
)

// POST /strings
func CreateStringHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Value == "" {
		http.Error(w, `"value" field cannot be empty`, http.StatusBadRequest)
		return
	}

	// Check if string already exists
	if _, _, _, err := database.GetStringByValue(req.Value); err == nil {
		http.Error(w, "string already exists", http.StatusConflict)
		return
	}

	// Analyze
	props, err := analysis.AnalyzeString(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	id := props["sha256_hash"].(string)

	// Insert into DB
	if err := database.InsertString(id, req.Value, props); err != nil {
		http.Error(w, "failed to save string", http.StatusInternalServerError)
		return
	}

	resp := models.StringReponsePayload{
		Id:         id,
		Value:      req.Value,
		Properties: props,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GET /strings
func GetAllStringsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if there are any filter parameters
	if len(r.URL.Query()) > 0 {
		handleFilteredStrings(w, r)
		return
	}

	results, err := database.GetAllStrings()
	if err != nil {
		http.Error(w, "failed to fetch strings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func handleFilteredStrings(w http.ResponseWriter, r *http.Request) {
	filters, err := filtering.ParseQueryParams(r)
	if err != nil {
		http.Error(w, "invalid filter parameters", http.StatusBadRequest)
		return
	}

	results, err := database.GetAllStrings()
	if err != nil {
		http.Error(w, "failed to fetch strings", http.StatusInternalServerError)
		return
	}

	filtered := filtering.ApplyFilters(results, filters)

	response := models.FilterPayload{
		Data:           filtered,
		Count:          len(filtered),
		FiltersApplied: filters,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func NaturalLanguageFilterHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	filters, err := nlp.ParseNaturalLanguageQuery(query)
	if err != nil {
		http.Error(w, "unable to parse natural language query", http.StatusBadRequest)
		return
	}

	results, err := database.GetAllStrings()
	if err != nil {
		http.Error(w, "failed to fetch strings", http.StatusInternalServerError)
		return
	}

	filtered := filtering.ApplyFilters(results, filters)

	response := models.NLPFilterPayload{
		Data:  filtered, // Now correctly using []models.StringReponsePayload
		Count: len(filtered),
		InterpretedQuery: models.InterpretedQuery{
			Original:      query,
			ParsedFilters: filters,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /strings/{string_value}
func GetStringHandler(w http.ResponseWriter, r *http.Request) {
	value := strings.TrimPrefix(r.URL.Path, "/strings/")
	if value == "" {
		http.Error(w, "string value not provided", http.StatusBadRequest)
		return
	}

	id, val, propsJSON, err := database.GetStringByValue(value)
	if err != nil {
		http.Error(w, "string not found", http.StatusNotFound)
		return
	}

	var props map[string]any
	if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
		http.Error(w, "failed to decode properties", http.StatusInternalServerError)
		return
	}

	resp := models.StringReponsePayload{
		Id:         id,
		Value:      val,
		Properties: props,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DELETE /strings/{string_value}
func DeleteStringHandler(w http.ResponseWriter, r *http.Request) {
	value := strings.TrimPrefix(r.URL.Path, "/strings/")
	if value == "" {
		http.Error(w, "string value not provided", http.StatusBadRequest)
		return
	}

	if err := database.DeleteString(value); err != nil {
		http.Error(w, "string not found or failed to delete", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
