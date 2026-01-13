package security

import (
	"encoding/hex"

	"golang.org/x/crypto/blake2b"
)

func PublicKeyToSuiAddress(pub []byte) string {
	scheme := byte(0x00)
	data := append([]byte{scheme}, pub...)

	hash := blake2b.Sum256(data)
	return "0x" + hex.EncodeToString(hash[:])
}
