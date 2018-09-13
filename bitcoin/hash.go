package bitcoin

import (
	"crypto/sha256"

	".."
)

var BitcoinHasher crypto.Hasher = btchash{}

type btchash struct{}

func (btchash) Hash(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash[:])
	return hash2[:], nil
}