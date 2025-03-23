package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(passHash), nil
}

func MatchPassword(dbPass, reqPass string) error {
	fmt.Printf("db %s, req %s\n", dbPass, reqPass)
	return bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(reqPass))
}
