package bitcoin

import (
	"bytes"
	"encoding/hex"

	"github.com/alivanz/go-utils"
)

type writer struct {
	version  uint32
	sequence uint32
	locktime uint32
	pubkey   []byte
}

func NewTxWriterV1(pubkey []byte) TxWriter {
	return &writer{1, 0xffffffff, 0, pubkey}
}

func (x *writer) TxWrite(txins []TxIn, txouts []TxOut) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := utils.NewBinaryWriter(buf)
	// version
	w.WriteUint32(x.version)
	// number of txin
	w.WriteCompactUint(uint64(len(txins)))
	for _, txin := range txins {
		// hash
		hash, _ := hex.DecodeString(txin.TxHash())
		w.Write(utils.ReverseBytes(hash))
		// output index
		w.WriteUint32(uint32(txin.OutputIndex()))
		// script
		w.WriteCompact(txin.Script())
		// sequence
		w.WriteUint32(x.sequence)
	}
	// number of txout
	w.WriteCompactUint(uint64(len(txouts)))
	for _, txout := range txouts {
		// amount to send
		w.WriteUint64(uint64(txout.Amount()))
		// script
		w.WriteCompact(txout.Script())
	}
	// locktime
	w.WriteUint32(x.locktime)
	// sig hash code
	w.WriteUint32(1)
	return buf.Bytes(), nil
}
