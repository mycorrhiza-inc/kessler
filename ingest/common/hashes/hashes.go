package hashes

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/blake2b"
)

type KesslerHash [32]byte

func HashFromBytes(b []byte) KesslerHash {
	result := blake2b.Sum256(b)
	return result
}

func (hash KesslerHash) String() string {
	return base64.URLEncoding.EncodeToString(hash[:])
}

func FromString(s string) (KesslerHash, error) {
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
