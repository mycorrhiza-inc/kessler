package hashes

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/blake2b"
)

// KesslerHash represents a base64 encoded BLAKE2b hash
// @Description A base64url-encoded BLAKE2b-256 hash
// @Schema {"type": "string", "example": "_EYNhTcsAPjIT3iNNvTnY5KFC1wm61Mki_uBcb3yKv2zDncVYfdI6c_7tH_PAAS8IlhNaapBg21fwT4Z7Ttxig=="}
type KesslerHash [32]byte

func (hash KesslerHash) String() string {
	return base64.URLEncoding.EncodeToString(hash[:])
}

func HashFromString(s string) (KesslerHash, error) {
	decoded, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return KesslerHash{}, err
	}
	if len(decoded) != 32 {
		return KesslerHash{}, fmt.Errorf("decoded base64 string length %d != 32", len(decoded))
	}
	var result KesslerHash
	copy(result[:], decoded)
	return result, nil
}

func (hash KesslerHash) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, hash.String())), nil
}

func (hash *KesslerHash) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}
	result, err := HashFromString(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*hash = result
	return nil
}

func HashFromBytes(b []byte) KesslerHash {
	result := blake2b.Sum256(b)
	return result
}

func HashFromFile(filePath string) (KesslerHash, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return KesslerHash{}, err
	}
	defer file.Close()
	var key []byte
	hash, err := blake2b.New256(key)
	if _, err := io.Copy(hash, file); err != nil {
		return KesslerHash{}, err
	}
	var result KesslerHash
	copy(result[:], hash.Sum(nil))
	return result, nil
}

func TestExpectedHash(input []byte, expectedHex string) error {
	expectedBytes, err := hex.DecodeString(expectedHex)
	if err != nil {
		return fmt.Errorf("Failed to decode hex string: %v", err)
	}

	var expectedHash KesslerHash
	copy(expectedHash[:], expectedBytes)

	computedHash := HashFromBytes(input)
	if computedHash != expectedHash {
		return fmt.Errorf("Hash mismatch\nExpected: %x\nGot:      %x", expectedHash, computedHash)
	}
	return nil
}
