package nlp

import (
	"strconv"
	"strings"
)

func ParseNaturalLanguageQuery(query string) (map[string]any, error) {
	filters := make(map[string]any)
	query = strings.ToLower(query)

	// Check for palindrome
	if strings.Contains(query, "palindrom") {
		filters["is_palindrome"] = true
	}

	// Check for word count
	if strings.Contains(query, "single word") {
		filters["word_count"] = 1
	}

	// Check for length constraints
	if strings.Contains(query, "longer than") {
		parts := strings.Split(query, "longer than")
		if len(parts) > 1 {
			numStr := strings.Split(strings.TrimSpace(parts[1]), " ")[0]
			if num, err := strconv.Atoi(numStr); err == nil {
				filters["min_length"] = num + 1
			}
		}
	}

	// Check for character contains
	for _, char := range []string{"a", "e", "i", "o", "u", "z"} {
		if strings.Contains(query, "contain "+char) ||
			strings.Contains(query, "contains "+char) {
			filters["contains_character"] = char
			break
		}
	}

	return filters, nil
}
