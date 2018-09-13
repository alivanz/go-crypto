package base58

import (
	"encoding/hex"
	"testing"
)

func TestEncode(t *testing.T) {
	data, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff")
	if Encode(data) != "1cWB5HCBdLjAuqGGReWE3R3CguuwSjw6RHn39s2yuDRTS5NsBgNiFpWgAnEx6VQi8csexkgYw3mdYrMHr8x9i7aEwP8kZ7vccXWqKDvGv3u1GxFKPuAkn8JCPPGDMf3vMMnbzm6Nh9zh1gcNsMvH3ZNLmP5fSG6DGbbi2tuwMWPthr4boWwCxf7ewSgNQeacyozhKDDQQ1qL5fQFUW52QKUZDZ5fw3KXNQJMcNTcaB723LchjeKun7MuGW5qyCBZYzA1KjofN1gYBV3NqyhQJ3Ns746GNuf9N2pQPmHz4xpnSrrfCvy6TVVz5d4PdrjeshsWQwpZsZGzvbdAdN8MKV5QsBDY" {
		t.Fail()
	}
}

func TestDecode(t *testing.T) {
	ref, _ := hex.DecodeString("73696d706c792061206c6f6e6720737472696e67")
	decoded, _ := Decode("2cFupjhnEsSn59qHXstmK2ffpLv2")
	if hex.EncodeToString(decoded) != hex.EncodeToString(ref) {
		t.Fail()
	}
}

func TestPrivWIF(t *testing.T) {
	ref, _ := hex.DecodeString("0ecd20654c2e2be708495853e8da35c664247040c00bd10b9b13e5e86e6a808d")
	decoded, _ := Decode("5HvofFG7K1e2aeWESm5pbCzRHtCSiZNbfLYXBvxyA57DhKHV4U3")
	// ignore wallet type and checksum
	wallethex := decoded[1:33]
	if hex.EncodeToString(wallethex) != hex.EncodeToString(ref) {
		t.Log(hex.EncodeToString(decoded))
		t.Fail()
	}
}
