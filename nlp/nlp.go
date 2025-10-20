package nlp

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseNaturalLanguageQuery converts natural language queries into filter parameters
// Supports patterns like:
// - "all single word palindromic strings" -> word_count=1, is_palindrome=true
// - "strings longer than 10 characters" -> min_length=11
// - "strings shorter than 5 characters" -> max_length=4
// - "palindromic strings that contain the letter a" -> is_palindrome=true, contains_character=a
// - "strings containing the letter z" -> contains_character=z
// - "exactly 3 word strings" -> word_count=3
func ParseNaturalLanguageQuery(query string) (map[string]any, error) {
	filters := make(map[string]any)
	queryLower := strings.ToLower(query)

	// Check for palindrome keywords
	if strings.Contains(queryLower, "palindrom") {
		filters["is_palindrome"] = true
	}

	// Check for word count patterns
	parseWordCount(queryLower, filters)

	// Check for length constraints
	parseLengthConstraints(queryLower, filters)

	// Check for character contains
	parseCharacterContains(queryLower, filters)

	return filters, nil
}

// parseWordCount handles patterns like:
// - "single word" -> 1
// - "two word" or "2 word" -> 2
// - "exactly 3 words" -> 3
func parseWordCount(query string, filters map[string]any) {
	// Pattern: "single word"
	if strings.Contains(query, "single word") {
		filters["word_count"] = 1
		return
	}

	// Pattern: "two word", "three word", etc.
	wordNumbers := map[string]int{
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
		"ten":   10,
	}

	for word, num := range wordNumbers {
		if strings.Contains(query, word+" word") {
			filters["word_count"] = num
			return
		}
	}

	// Pattern: "exactly 3 words" or "3 words"
	re := regexp.MustCompile(`(?:exactly\s+)?(\d+)\s+words?`)
	if match := re.FindStringSubmatch(query); match != nil {
		if num, err := strconv.Atoi(match[1]); err == nil {
			filters["word_count"] = num
		}
	}
}

// parseLengthConstraints handles patterns like:
// - "longer than 10 characters" -> min_length=11
// - "shorter than 5 characters" -> max_length=4
// - "at least 8 characters" -> min_length=8
// - "at most 20 characters" -> max_length=20
func parseLengthConstraints(query string, filters map[string]any) {
	// Pattern: "longer than X" or "longer than X characters"
	re := regexp.MustCompile(`longer than\s+(\d+)`)
	if match := re.FindStringSubmatch(query); match != nil {
		if num, err := strconv.Atoi(match[1]); err == nil {
			filters["min_length"] = num + 1
		}
		return
	}

	// Pattern: "shorter than X" or "shorter than X characters"
	re = regexp.MustCompile(`shorter than\s+(\d+)`)
	if match := re.FindStringSubmatch(query); match != nil {
		if num, err := strconv.Atoi(match[1]); err == nil {
			filters["max_length"] = num - 1
		}
		return
	}

	// Pattern: "at least X characters"
	re = regexp.MustCompile(`at least\s+(\d+)`)
	if match := re.FindStringSubmatch(query); match != nil {
		if num, err := strconv.Atoi(match[1]); err == nil {
			filters["min_length"] = num
		}
		return
	}

	// Pattern: "at most X characters"
	re = regexp.MustCompile(`at most\s+(\d+)`)
	if match := re.FindStringSubmatch(query); match != nil {
		if num, err := strconv.Atoi(match[1]); err == nil {
			filters["max_length"] = num
		}
		return
	}

	// Pattern: "exactly X characters"
	re = regexp.MustCompile(`exactly\s+(\d+)`)
	if match := re.FindStringSubmatch(query); match != nil {
		if num, err := strconv.Atoi(match[1]); err == nil {
			filters["min_length"] = num
			filters["max_length"] = num
		}
	}
}

// parseCharacterContains handles patterns like:
// - "containing the letter a" -> contains_character=a
// - "contains the letter z" -> contains_character=z
// - "that contain the first vowel" -> contains_character=a
// - "with character x" -> contains_character=x
func parseCharacterContains(query string, filters map[string]any) {
	// Pattern: "contain(s) the letter X" or "contain(s) the character X"
	re := regexp.MustCompile(`contains?\s+the\s+(?:letter|character)\s+([a-z])`)
	if match := re.FindStringSubmatch(query); match != nil {
		filters["contains_character"] = match[1]
		return
	}

	// Pattern: "with character X"
	re = regexp.MustCompile(`with\s+(?:character|letter)\s+([a-z])`)
	if match := re.FindStringSubmatch(query); match != nil {
		filters["contains_character"] = match[1]
		return
	}

	// Pattern: "first vowel" - default to 'a'
	if strings.Contains(query, "first vowel") {
		filters["contains_character"] = "a"
		return
	}

	// Pattern: "vowel" without specifier - could be improved with more context
	// For now, we'll check for specific vowel mentions
	vowels := []string{"a", "e", "i", "o", "u"}
	for _, vowel := range vowels {
		if strings.Contains(query, "contain "+vowel) ||
			strings.Contains(query, "contains "+vowel) ||
			strings.Contains(query, "letter "+vowel) ||
			strings.Contains(query, "character "+vowel) {
			filters["contains_character"] = vowel
			return
		}
	}
}
