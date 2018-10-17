package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"

	utils "github.com/alivanz/go-utils"
)

type TXBuilder struct {
	Version  uint32
	Sequence uint32
	Locktime uint32
	txins    []TxIn
	txouts   []TxOut
}

type customtxin struct {
	txhash   string
	outindex int
	script   []byte
}

func CustomTxIn(txhash string, outindex int, script []byte) TxIn {
	return &customtxin{txhash, outindex, script}
}
func (tx *customtxin) TxHash() string {
	return tx.txhash
}
func (tx *customtxin) OutputIndex() int {
	return tx.outindex
}
func (tx *customtxin) Script() []byte {
	return tx.script
}

func NewTXBuilder(version uint32) *TXBuilder {
	return &TXBuilder{version, 0xffffffff, 0, nil, nil}
}
func (builder *TXBuilder) AddTxIn(tx TxIn) {
	builder.txins = append(builder.txins, tx)
}
func (builder *TXBuilder) AddTxOut(tx TxOut) {
	builder.txouts = append(builder.txouts, tx)
}

func (builder *TXBuilder) RawTransaction() []byte {
	buf := bytes.NewBuffer(nil)
	w := utils.NewBinaryWriter(buf)
	// version
	w.WriteUint32(builder.Version)
	// txin
	w.WriteCompactUint(uint64(len(builder.txins)))
	for _, tx := range builder.txins {
		// w.Write(txin)
		hash, _ := hex.DecodeString(tx.TxHash())
		w.Write(utils.ReverseBytes(hash))
		// output index
		w.WriteUint32(uint32(tx.OutputIndex()))
		// script
		w.WriteCompact(tx.Script())
		// sequence
		w.WriteUint32(builder.Sequence)
	}
	// number of txout
	w.WriteCompactUint(uint64(len(builder.txouts)))
	for _, tx := range builder.txouts {
		// w.Write(txout)
		w.WriteUint64(uint64(tx.Amount()))
		// script
		w.WriteCompact(tx.Script())
	}
	// locktime
	w.WriteUint32(builder.Locktime)
	return buf.Bytes()
}

func (builder *TXBuilder) unsignedTransaction(i int) []byte {
	x := NewTXBuilder(builder.Version)
	for j, tx := range builder.txins {
		if i == j {
			x.AddTxIn(tx)
		} else {
			x.AddTxIn(CustomTxIn(tx.TxHash(), tx.OutputIndex(), nil))
		}
	}
	for _, tx := range builder.txouts {
		x.AddTxOut(tx)
	}
	return x.RawTransaction()
}

func (builder *TXBuilder) MessageHashes() [][]byte {
	list := make([][]byte, len(builder.txins))
	for i := 0; i < len(builder.txins); i++ {
		msg := append(builder.unsignedTransaction(i), 0x01, 0x00, 0x00, 0x00)
		hash := BitcoinHasher.Hash(msg)
		list[i] = hash
	}
	return list
}

func (builder *TXBuilder) SignedTransaction(signscripts [][]byte) ([]byte, error) {
	if len(builder.txins) != len(signscripts) {
		return nil, fmt.Errorf("number of script doesnt match")
	}
	x := NewTXBuilder(builder.Version)
	for i, tx := range builder.txins {
		x.AddTxIn(CustomTxIn(tx.TxHash(), tx.OutputIndex(), signscripts[i]))
	}
	for _, tx := range builder.txouts {
		x.AddTxOut(tx)
	}
	return x.RawTransaction(), nil
}
