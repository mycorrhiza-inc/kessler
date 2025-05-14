package hashes_test

import (
	"crypto/rand"
	"kessler/pkg/hashes"
	"testing"
)

func TestKesslerHashRoundTrip(t *testing.T) {
	for i := 0; i < 1000; i++ {
		// Generate random 32-byte KesslerHash
		// Or better, fill with random bytes:
		original := hashes.KesslerHash{}
		rand.Read(original[:])

		// Convert to base64 string
		s := original.String()

		// Convert back from string
		decoded, err := hashes.HashFromString(s)
		if err != nil {
			t.Errorf("Error decoding string on iteration %d: %v", i, err)
			continue
		}

		// Verify decoded hash matches original
		if original != decoded {
			t.Errorf("Decoded hash does not match original on iteration %d", i)
		}
	}
}

func TestKesslerHashValidity(t *testing.T) {
	input := []byte("The quick brown fox jumped over the lazy dog")
	expectedHex := "cd1c3b120f8d0af28a9b6b1c43da5aba4be633ac0a303719f6dfa5ee1890f28d"
	err := hashes.TestExpectedHash(input, expectedHex)
	if err != nil {
		t.Fatalf("Hash did not match for input %v: %v", input, err)
	}
	input = []byte("the mitochondria is the powerhouse of a cell")
	expectedHex = "821435d2a2b379ad2e4bb11c41c0b2ec2cf2135f09b0afa740d5efc2818778f7"
	err = hashes.TestExpectedHash(input, expectedHex)
	if err != nil {
		t.Fatalf("Hash did not match for input %v: %v", input, err)
	}
}
