package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateVerificationCode() int {
	code, _ := rand.Int(rand.Reader, big.NewInt(10)) // random number [0-9]
	return int(code.Int64())
}
