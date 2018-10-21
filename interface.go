package crypto

import (
	"crypto/ecdsa"
	"math/big"
)

// Transaction maker, single recipient
type TxMaker interface {
	TxMake(from string, to string) ([]byte, error)
}

// Transaction Hasher
type Hasher interface {
	Hash(data []byte) []byte
}

type Wallet interface {
	Signer
	Verifier
	PubKey() (*ecdsa.PublicKey, error)
}

// Messge Signer
type Signer interface {
	Sign(hash []byte) (*big.Int, *big.Int, error)
}

// Message Signature Verifier
type Verifier interface {
	Verify(hash []byte, r, s *big.Int) bool
}

// Smallest unit
type Unit interface {
	String() string
}
