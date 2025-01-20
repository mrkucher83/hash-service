package cryptographer

import (
	"fmt"
	"golang.org/x/crypto/sha3"
)

func Hash(input []string) []string {
	var hashes []string

	for _, str := range input {
		hash := sha3.Sum256([]byte(str))
		hashes = append(hashes, fmt.Sprintf("%x", hash))
	}
	return hashes
}
