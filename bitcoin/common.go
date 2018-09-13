package bitcoin

type simpletxout struct {
	address string
	amount  Satoshi
}

func (tx *simpletxout) Address() string {
	return tx.address
}
func (tx *simpletxout) Amount() Satoshi {
	return tx.amount
}

// New TxOut
func NewTxOut(address string, amount Satoshi) TxOut {
	return &simpletxout{address, amount}
}
