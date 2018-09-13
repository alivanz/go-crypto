package crypto

// Transaction maker, single recipient
type TxMaker interface {
	TxMake(from string, to string) ([]byte, error)
}

// Transaction Hasher
type Hasher interface {
	Hash(data []byte) ([]byte, error)
}

// Transaction Signer
type Signer interface {
	Sign(hash []byte) ([]byte, error)
}

// Smallest unit
type Unit interface {
	String() string
}
