package pkg

import (
	"crypto/rand"
	"math/big"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func BigSerial() *big.Int {
	serial, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
	return serial
}
