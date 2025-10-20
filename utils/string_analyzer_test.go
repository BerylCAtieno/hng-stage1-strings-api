package utils

import (
	"reflect"
	"testing"
)

// TestIsPalindrome checks various palindrome cases.
func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"racecar", true},
		{"A man, a plan, a canal: Panama", true},
		{"No lemon, no melon", true},
		{"hello", false},
		{"😊abccba😊", true},
		{"12321", true},
		{"12345", false},
	}

	for _, test := range tests {
		result := IsPalindrome(test.input)
		if result != test.expected {
			t.Errorf("IsPalindrome(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

// TestLength verifies rune-based length counting.
func TestLength(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"😊😊", 2}, // Unicode characters
		{"", 0},
		{"GoLang", 6},
	}

	for _, test := range tests {
		result := Length(test.input)
		if result != test.expected {
			t.Errorf("Length(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}

// TestNumberOfWords checks word counting with different spacing.
func TestNumberOfWords(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello world", 2},
		{" multiple   spaces here ", 3},
		{"", 0},
		{"one-word", 1},
		{"line\nbreak test", 3},
	}

	for _, test := range tests {
		result := NumberOfWords(test.input)
		if result != test.expected {
			t.Errorf("NumberOfWords(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}

// TestUniqueCharacters checks case-insensitive unique character counting.
func TestUniqueCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 4}, // h,e,l,o
		{"AaBb", 2},  // a,b (case-insensitive)
		{"", 0},
		{"Go😊Go", 3}, // g,o,😊
	}

	for _, test := range tests {
		result := UniqueCharacters(test.input)
		if result != test.expected {
			t.Errorf("UniqueCharacters(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}

// TestCharFrequency verifies character frequency counting (ignoring spaces).
func TestCharFrequency(t *testing.T) {
	tests := []struct {
		input    string
		expected map[rune]int
	}{
		{"hello world", map[rune]int{
			'h': 1, 'e': 1, 'l': 3, 'o': 2, 'w': 1, 'r': 1, 'd': 1,
		}},
		{"AaBb", map[rune]int{
			'a': 2, 'b': 2,
		}},
		{"😊😊 Go", map[rune]int{
			'😊': 2, 'g': 1, 'o': 1,
		}},
		{"   ", map[rune]int{}}, // only spaces
	}

	for _, test := range tests {
		result := CharFrequency(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("CharFrequency(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
