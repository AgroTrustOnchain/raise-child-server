package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"raise-child/constants/env"
)

func Encrypt(text string) string {
	var key string = os.Getenv(env.ENCRYPT_SECRET_KEY)
	block, _ := aes.NewCipher([]byte(key))
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	var ciphertext = gcm.Seal(nonce, nonce, []byte(text), nil)
	return hex.EncodeToString(ciphertext)
}

func Decrypt(text string) string {
	var key string = os.Getenv(env.ENCRYPT_SECRET_KEY)
	data, _ := hex.DecodeString(text)
	block, _ := aes.NewCipher([]byte(key))
	gcm, _ := cipher.NewGCM(block)

	var nonceSize int = gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return string(plaintext)
}
