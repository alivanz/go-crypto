package bitcoin

// Address information
type AddressInfo interface {
	Address() string
	Balance() Satoshi
	UnconfirmedBalance() Satoshi
}

// Supply transaction
type TxSupplier interface {
	AllTxSupply() ([]Tx, error)
	PastTxSupplier
	UnspentTxSupplier
}

// Supply past transaction
type PastTxSupplier interface {
	PastTxSupply() ([]Tx, error)
}

// Supply unspent transaction
type UnspentTxSupplier interface {
	UnspentTxSupply() ([]TxIn, error)
}

// Transaction filter
type TxFilter interface {
	TxFilter(change_script []byte, unspents []TxIn, recipient []TxOut, fee Satoshi) ([]TxIn, []TxOut, error)
}

// Raw Tx Writer (bitcoin related)
type TxWriter interface {
	TxWrite(txins []TxIn, txouts []TxOut) ([]byte, error)
}

// Transaction
type Tx interface {
	TxHash() string
	Confirmation() int
	BlockHeight() int
	Spent() bool
}
type TxFull interface {
	Tx
	Inputs() []TxIn
	Outputs() []TxOut
}

// TxIn (bitcoin related)
type TxIn interface {
	TxHash() string
	OutputIndex() int
	Amount() Satoshi
	Script() []byte
}

// TxOut (bitcoin related)
type TxOut interface {
	Amount() Satoshi
	Script() []byte
}

// Broadcast raw transaction
type Broadcaster interface {
	Broadcast(rawtx []byte) error
}
