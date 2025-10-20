package filtering

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/BerylCAtieno/string-analyzer/models"
)

// ParseQueryParams extracts and validates filter parameters from URL query
func ParseQueryParams(r *http.Request) (map[string]any, error) {
	filters := make(map[string]any)
	query := r.URL.Query()

	// Parse is_palindrome
	if val := query.Get("is_palindrome"); val != "" {
		isPalindrome, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		filters["is_palindrome"] = isPalindrome
	}

	// Parse min_length
	if val := query.Get("min_length"); val != "" {
		minLength, err := strconv.Atoi(val)
		if err != nil || minLength < 0 {
			return nil, err
		}
		filters["min_length"] = minLength
	}

	// Parse max_length
	if val := query.Get("max_length"); val != "" {
		maxLength, err := strconv.Atoi(val)
		if err != nil || maxLength < 0 {
			return nil, err
		}
		filters["max_length"] = maxLength
	}

	// Parse word_count
	if val := query.Get("word_count"); val != "" {
		wordCount, err := strconv.Atoi(val)
		if err != nil || wordCount < 0 {
			return nil, err
		}
		filters["word_count"] = wordCount
	}

	// Parse contains_character
	if val := query.Get("contains_character"); val != "" {
		if len(val) != 1 {
			return nil, errors.New("contains_character must be a single character")
		}
		filters["contains_character"] = strings.ToLower(val)
	}

	return filters, nil
}

// ApplyFilters filters a slice of strings based on given criteria
// ApplyFilters filters a slice of StringReponsePayload based on given criteria
func ApplyFilters(data []models.StringReponsePayload, filters map[string]any) []models.StringReponsePayload {
	var filtered []models.StringReponsePayload

	for _, item := range data {
		if matchesFilters(item, filters) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func matchesFilters(item models.StringReponsePayload, filters map[string]any) bool {
	props := item.Properties

	// Check is_palindrome
	if val, ok := filters["is_palindrome"]; ok {
		if props["is_palindrome"] != val {
			return false
		}
	}

	// Check min_length
	if val, ok := filters["min_length"]; ok {
		if props["length"].(int) < val.(int) {
			return false
		}
	}

	// Check max_length
	if val, ok := filters["max_length"]; ok {
		if props["length"].(int) > val.(int) {
			return false
		}
	}

	// Check word_count
	if val, ok := filters["word_count"]; ok {
		if props["word_count"] != val {
			return false
		}
	}

	// Check contains_character
	if char, ok := filters["contains_character"]; ok {
		freqMap := props["character_frequency_map"].(map[string]int)
		if freqMap[char.(string)] == 0 {
			return false
		}
	}

	return true
}
