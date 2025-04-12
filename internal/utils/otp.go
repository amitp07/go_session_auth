package utils

import (
	"fmt"
	"math/rand"
)

func GenerateOtp() string {
	return fmt.Sprintf("%06d", rand.Intn(1_000_000))
}
