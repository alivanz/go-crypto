package bitcoin

import "encoding/hex"

type simpletxout struct {
	script []byte
	amount Satoshi
}
type simpletxin struct {
	txhash string
	index  int
	amount Satoshi
	script []byte
}

func (tx *simpletxout) Script() []byte {
	return tx.script
}
func (tx *simpletxout) Amount() Satoshi {
	return tx.amount
}
func (tx *simpletxin) TxHash() string {
	return tx.txhash
}
func (tx *simpletxin) Amount() Satoshi {
	return tx.amount
}
func (tx *simpletxin) OutputIndex() int {
	return tx.index
}
func (tx *simpletxin) Script() []byte {
	return tx.script
}

// New TxIn
func NewTxIn(txhash string, index int, amount Satoshi, script string) TxInExtended {
	x, _ := hex.DecodeString(script)
	return &simpletxin{txhash, index, amount, x}
}

// New TxOut
func NewTxOut(script []byte, amount Satoshi) TxOut {
	return &simpletxout{script, amount}
}
