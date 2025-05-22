package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateVerificationCode() string {
	code := ""
	for i := 0; i < 5; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10)) // random number [0-9]

		code += n.String()
	}
	return code
}
