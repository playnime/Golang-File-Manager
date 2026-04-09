package myrsa

import "math/big"

func Decrypt(AES_с *big.Int, private *big.Int, n *big.Int) *big.Int {
	return new(big.Int).Exp(AES_с, private, n)
}
