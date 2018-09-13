package bitcoin

import "fmt"

type noTxFilter struct{}
type minTxFilter struct{}

var (
	NoTxFilter       TxFilter = noTxFilter{}
	MinTxFilter      TxFilter = minTxFilter{}
	UnsufficientFund          = fmt.Errorf("Unsufficient fund")
)

func (noTxFilter) TxFilter(change_address string, unspents []TxIn, recipient []TxOut, fee Satoshi) ([]TxIn, []TxOut, error) {
	total_in := Satoshi(0)
	total_out_fee := fee
	for _, tx := range unspents {
		total_in = total_in + tx.Amount()
	}
	for _, tx := range recipient {
		total_out_fee = total_out_fee + tx.Amount()
	}
	if total_in < total_out_fee {
		return nil, nil, UnsufficientFund
	}
	change := total_in - total_out_fee
	return unspents, append(recipient, NewTxOut(change_address, change)), nil
}

func (minTxFilter) TxFilter(change_address string, unspents []TxIn, recipient []TxOut, fee Satoshi) ([]TxIn, []TxOut, error) {
	total_out_fee := fee
	for _, tx := range recipient {
		total_out_fee = total_out_fee + tx.Amount()
	}
	txins := make([]TxIn, 0)
	total_in := Satoshi(0)
	for _, tx := range unspents {
		if total_in < total_out_fee {
			txins = append(txins, tx)
			total_in = total_in + tx.Amount()
		}
	}
	if total_in < total_out_fee {
		return nil, nil, UnsufficientFund
	}
	change := total_in - total_out_fee
	return txins, append(recipient, NewTxOut(change_address, change)), nil
}
