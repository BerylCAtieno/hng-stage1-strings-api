package analysis

import (
	"errors"

	"github.com/BerylCAtieno/string-analyzer/utils"
)

// Properties holds computed attributes of a string
type Properties map[string]any

// AnalyzeString computes various attributes of an input string
// Returns an error if the string is empty
func AnalyzeString(s string) (Properties, error) {

	if len([]rune(s)) == 0 {
		return nil, errors.New("string cannot be empty")
	}

	properties := Properties{
		"length":                  utils.Length(s),
		"is_palindrome":           utils.IsPalindrome(s),
		"unique_characters":       utils.UniqueCharacters(s),
		"word_count":              utils.NumberOfWords(s),
		"sha256_hash":             utils.GenerateId(s),
		"character_frequency_map": utils.CharFrequency(s),
	}

	return properties, nil

}
