package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRadomId(len int) (string, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
