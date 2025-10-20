package utils

import (
	"encoding/hex"
	"testing"
)

// TestGenerateID_SameInputSameHash ensures determinism:
// same input → same hash every time.
func TestGenerateID_SameInputSameHash(t *testing.T) {
	input := "hello world"
	id1 := GenerateId(input)
	id2 := GenerateId(input)

	if id1 != id2 {
		t.Errorf("expected same hash for identical inputs, got %s and %s", id1, id2)
	}

	// SHA-256 hex should always be 64 characters long.
	if len(id1) != 64 {
		t.Errorf("expected hash length 64, got %d", len(id1))
	}

	// Verify it’s valid hexadecimal.
	if _, err := hex.DecodeString(id1); err != nil {
		t.Errorf("hash is not valid hexadecimal: %v", err)
	}
}

// TestGenerateID_DifferentInputsDifferentHash ensures distinct inputs
// produce different hashes (extremely likely).
func TestGenerateID_DifferentInputsDifferentHash(t *testing.T) {
	id1 := GenerateId("apple")
	id2 := GenerateId("banana")

	if id1 == id2 {
		t.Errorf("expected different hashes for different inputs, got same: %s", id1)
	}
}

// TestGenerateID_EmptyString ensures that even empty input yields a hash.
func TestGenerateID_EmptyString(t *testing.T) {
	id := GenerateId("")
	if id == "" {
		t.Error("expected non-empty hash for empty string input")
	}
	if len(id) != 64 {
		t.Errorf("expected 64-character hash, got %d", len(id))
	}
}
