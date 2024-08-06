package ust

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func NewRandomID() int {
	num, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		fmt.Println("Error al generar un número aleatorio", err)
		return 0
	}
	return int(num.Int64())
}

func NewPersonaID() int {
	num, err := rand.Int(rand.Reader, big.NewInt(9999))
	if err != nil {
		fmt.Println("Error al generar un número aleatorio", err)
		return 0
	}
	return int(num.Int64())
}
