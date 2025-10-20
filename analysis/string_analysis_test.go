package analysis

import (
	"reflect"
	"testing"
)

// TestAnalyzeString_ValidInput tests that AnalyzeString returns correct values for a sample input.
func TestAnalyzeString_ValidInput(t *testing.T) {
	input := "hello world"

	props, err := AnalyzeString(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// --- Basic existence checks ---
	requiredKeys := []string{
		"length",
		"is_palindrome",
		"unique_characters",
		"word_count",
		"sha256_hash",
		"character_frequency_map",
	}

	for _, key := range requiredKeys {
		if _, ok := props[key]; !ok {
			t.Errorf("expected key %q in properties map", key)
		}
	}

	// --- Type checks ---
	if _, ok := props["length"].(int); !ok {
		t.Errorf("expected 'length' to be int")
	}
	if _, ok := props["is_palindrome"].(bool); !ok {
		t.Errorf("expected 'is_palindrome' to be bool")
	}
	if _, ok := props["unique_characters"].(int); !ok {
		t.Errorf("expected 'unique_characters' to be int")
	}
	if _, ok := props["word_count"].(int); !ok {
		t.Errorf("expected 'word_count' to be int")
	}
	if _, ok := props["sha256_hash"].(string); !ok {
		t.Errorf("expected 'sha256_hash' to be string")
	}
	if _, ok := props["character_frequency_map"].(map[rune]int); !ok {
		t.Errorf("expected 'character_frequency_map' to be map[rune]int")
	}

	// --- Value checks ---
	if props["is_palindrome"].(bool) {
		t.Errorf("expected 'hello world' not to be palindrome")
	}
	if props["word_count"].(int) != 2 {
		t.Errorf("expected word count 2, got %v", props["word_count"])
	}
	if props["unique_characters"].(int) <= 0 {
		t.Errorf("expected positive unique character count, got %v", props["unique_characters"])
	}

	// --- Frequency map check ---
	freq, _ := props["character_frequency_map"].(map[rune]int)
	expectedFreq := map[rune]int{
		'h': 1, 'e': 1, 'l': 3, 'o': 2, 'w': 1, 'r': 1, 'd': 1,
	}
	if !reflect.DeepEqual(freq, expectedFreq) {
		t.Errorf("unexpected frequency map.\nGot: %#v\nWant: %#v", freq, expectedFreq)
	}
}

// TestAnalyzeString_EmptyInput ensures empty strings return an error.
func TestAnalyzeString_EmptyInput(t *testing.T) {
	props, err := AnalyzeString("")
	if err == nil {
		t.Fatalf("expected error for empty input, got nil")
	}
	if props != nil {
		t.Fatalf("expected nil properties map for empty input")
	}
}
