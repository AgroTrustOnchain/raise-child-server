package util

import (
	"crypto/rand"
	"math/big"

	"github.com/google/uuid"
)

func GenerateNonce() string {
	return uuid.NewString()
}

func GenerateSalt() string {
	// 1. Create a 16-byte buffer/array
	var bytes = make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		var id = uuid.New()
		return new(big.Int).SetBytes(id[:]).String()
	}

	return new(big.Int).SetBytes(bytes).String()
}

func GenerateString() string {
	return uuid.NewString()
}

func GenerateId() string {
	if res, err := uuid.NewV7(); err == nil {
		return res.String()
	}

	return uuid.NewString()
}
