package crypto

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
	PubKey() []byte
}

// Messge Signer
type Signer interface {
	Sign(hash []byte) ([]byte, error)
}

// Message Signature Verifier
type Verifier interface {
	Verify(hash []byte, sig []byte) bool
}

// Smallest unit
type Unit interface {
	String() string
}
