package hashes_test

import (
	"crypto/rand"
	"testing"
	"thaumaturgy/common/hashes"
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
