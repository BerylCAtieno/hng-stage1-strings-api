package utils

import (
	"strings"
	"unicode"
)

// IsPalindrome checks if a string is a palindrome.
func IsPalindrome(s string) bool {
	var cleaned []rune
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			cleaned = append(cleaned, unicode.ToLower(char))
		}
	}
	for i, j := 0, len(cleaned)-1; i < j; i, j = i+1, j-1 {
		if cleaned[i] != cleaned[j] {
			return false
		}
	}
	return true
}

// Length returns the number of characters (not bytes) in a string.
func Length(s string) int {
	return len([]rune(s))
}

// NumberOfWords counts the number of words in a string.
func NumberOfWords(s string) int {
	return len(strings.Fields(s))
}

// UniqueCharacters counts the number of unique characters in a string.
func UniqueCharacters(s string) int {
	unique := make(map[rune]struct{})
	for _, char := range s {
		unique[unicode.ToLower(char)] = struct{}{}
	}
	return len(unique)
}

// CharFrequency returns a map of character frequencies in a string.
func CharFrequency(s string) map[rune]int {
	frequency := make(map[rune]int)
	for _, char := range s {
		if unicode.IsSpace(char) {
			continue
		}
		frequency[unicode.ToLower(char)]++
	}
	return frequency
}
