package myrsa

import "math/big"

func Encrypt(AES *big.Int, public *big.Int, n *big.Int) *big.Int {
	return new(big.Int).Exp(AES, public, n)
}
