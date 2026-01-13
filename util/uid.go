package util

import "github.com/google/uuid"

func GenerateNonce() string {
	return uuid.NewString()
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
