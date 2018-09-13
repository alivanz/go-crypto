package base58

// #include <stdint.h>
// #include "libbase58.h"
import "C"
import (
	"errors"
	"unsafe"
)

var Chars = []byte("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")

func Encode(src []byte) string {
	b58 := make([]byte, len(src)*5)
	b58_len := C.size_t(len(b58))
	success := bool(C.b58enc(
		(*C.char)(unsafe.Pointer(&b58[0])),
		&b58_len,
		unsafe.Pointer(&src[0]),
		C.size_t(len(src)),
	))
	if !success {
		panic("Encode Error")
	}
	return string(b58[:int(b58_len)-1])
}
func Decode(src string) ([]byte, error) {
	bin := make([]byte, len(src))
	bin_len := C.size_t(len(bin))
	chars := C.CString(src)
	success := bool(C.b58tobin(
		unsafe.Pointer(&bin[0]),
		&bin_len,
		chars,
		C.size_t(len(src)),
	))
	if !success {
		return nil, errors.New("Decode Error")
	}
	return bin[len(bin)-int(bin_len):], nil
}
