package filtering

import (
	"errors"
	"fmt"
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
			return nil, fmt.Errorf("invalid is_palindrome value: %v", err)
		}
		filters["is_palindrome"] = isPalindrome
	}

	// Parse min_length
	var minLength int
	if val := query.Get("min_length"); val != "" {
		var err error
		minLength, err = strconv.Atoi(val)
		if err != nil || minLength < 0 {
			return nil, fmt.Errorf("invalid min_length value: must be a non-negative integer")
		}
		filters["min_length"] = minLength
	}

	// Parse max_length
	var maxLength int
	if val := query.Get("max_length"); val != "" {
		var err error
		maxLength, err = strconv.Atoi(val)
		if err != nil || maxLength < 0 {
			return nil, fmt.Errorf("invalid max_length value: must be a non-negative integer")
		}
		filters["max_length"] = maxLength
	}

	// Validate that min_length <= max_length
	if minVal, hasMin := filters["min_length"]; hasMin {
		if maxVal, hasMax := filters["max_length"]; hasMax {
			if minVal.(int) > maxVal.(int) {
				return nil, errors.New("min_length cannot be greater than max_length")
			}
		}
	}

	// Parse word_count
	if val := query.Get("word_count"); val != "" {
		wordCount, err := strconv.Atoi(val)
		if err != nil || wordCount < 0 {
			return nil, fmt.Errorf("invalid word_count value: must be a non-negative integer")
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

// ApplyFilters filters a slice of StringResponsePayload based on given criteria
func ApplyFilters(data []models.StringReponsePayload, filters map[string]any) []models.StringReponsePayload {
	var filtered []models.StringReponsePayload

	for _, item := range data {
		if matchesFilters(item, filters) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// matchesFilters checks if a single StringResponsePayload matches all applied filters
func matchesFilters(item models.StringReponsePayload, filters map[string]any) bool {
	props := item.Properties

	// Check is_palindrome filter
	if val, ok := filters["is_palindrome"]; ok {
		isPalindrome, err := safeGetBool(props, "is_palindrome")
		if err != nil {
			return false
		}
		if isPalindrome != val.(bool) {
			return false
		}
	}

	// Check min_length filter
	if val, ok := filters["min_length"]; ok {
		length, err := safeGetInt(props, "length")
		if err != nil {
			return false
		}
		if length < val.(int) {
			return false
		}
	}

	// Check max_length filter
	if val, ok := filters["max_length"]; ok {
		length, err := safeGetInt(props, "length")
		if err != nil {
			return false
		}
		if length > val.(int) {
			return false
		}
	}

	// Check word_count filter
	if val, ok := filters["word_count"]; ok {
		wordCount, err := safeGetInt(props, "word_count")
		if err != nil {
			return false
		}
		if wordCount != val.(int) {
			return false
		}
	}

	// Check contains_character filter
	if char, ok := filters["contains_character"]; ok {
		freqMap, err := safeGetCharFrequencyMap(props, "character_frequency_map")
		if err != nil {
			return false
		}
		charStr := char.(string)
		// Check if character exists and has count > 0
		if count, exists := freqMap[charStr]; !exists || count == 0 {
			return false
		}
	}

	return true
}

// safeGetInt safely retrieves an integer value from properties
func safeGetInt(props map[string]any, key string) (int, error) {
	val, ok := props[key]
	if !ok {
		return 0, fmt.Errorf("property '%s' not found", key)
	}

	intVal, ok := val.(int)
	if !ok {
		// Try float64 (common from JSON unmarshaling)
		if floatVal, ok := val.(float64); ok {
			return int(floatVal), nil
		}
		return 0, fmt.Errorf("property '%s' is not an integer", key)
	}

	return intVal, nil
}

// safeGetBool safely retrieves a boolean value from properties
func safeGetBool(props map[string]any, key string) (bool, error) {
	val, ok := props[key]
	if !ok {
		return false, fmt.Errorf("property '%s' not found", key)
	}

	boolVal, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("property '%s' is not a boolean", key)
	}

	return boolVal, nil
}

// safeGetCharFrequencyMap safely retrieves the character frequency map from properties
func safeGetCharFrequencyMap(props map[string]any, key string) (map[string]int, error) {
	val, ok := props[key]
	if !ok {
		return nil, fmt.Errorf("property '%s' not found", key)
	}

	// Handle map[string]interface{} (common from JSON unmarshaling)
	freqMapInterface, ok := val.(map[string]interface{})
	if ok {
		freqMap := make(map[string]int)
		for k, v := range freqMapInterface {
			// Convert each value to int
			switch v := v.(type) {
			case float64:
				freqMap[k] = int(v)
			case int:
				freqMap[k] = v
			default:
				return nil, fmt.Errorf("character frequency map contains non-numeric value for '%s'", k)
			}
		}
		return freqMap, nil
	}

	// Handle map[string]int directly
	freqMap, ok := val.(map[string]int)
	if !ok {
		return nil, fmt.Errorf("property '%s' is not a valid character frequency map", key)
	}

	return freqMap, nil
}
