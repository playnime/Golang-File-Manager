package myrsa

import (
	"crypto/rand"
	"math/big"
)

func randBigInt() (*big.Int, error) {
	byteLen := 512 / 8

	bytes := make([]byte, byteLen)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	bytes[0] |= 1 << 7
	bytes[byteLen-1] |= 1

	n := new(big.Int).SetBytes(bytes)
	return n, nil
}

func checkOnPrime(n *big.Int) bool {
	zero := big.NewInt(0)
	nbig := new(big.Int).Set(n)

	for i := int64(2); i < 1000; i++ {
		iBig := big.NewInt(i)

		if nbig.Cmp(iBig) == 0 {
			return true
		}

		rem := new(big.Int).Rem(nbig, iBig)
		if rem.Cmp(zero) == 0 {
			return false
		}
	}

	return nbig.ProbablyPrime(20)
}

func generatePQ() (*big.Int, *big.Int, *big.Int, *big.Int) {
	e := big.NewInt(65537)
	one := big.NewInt(1)

	for {
		p, err := randBigInt()
		if err != nil {
			return nil, nil, nil, nil
		}
		for !checkOnPrime(p) {
			p, err = randBigInt()
			if err != nil {
				return nil, nil, nil, nil
			}
		}

		q, err := randBigInt()
		if err != nil {
			return nil, nil, nil, nil
		}
		for !checkOnPrime(q) {
			q, err = randBigInt()
			if err != nil {
				return nil, nil, nil, nil
			}
		}

		if p.Cmp(q) == 0 {
			continue
		}

		n := new(big.Int).Mul(p, q)
		if n.BitLen() != 1024 {
			continue
		}

		pm1 := new(big.Int).Sub(p, one)
		qm1 := new(big.Int).Sub(q, one)
		phi := new(big.Int).Mul(pm1, qm1)

		gcd := new(big.Int).GCD(nil, nil, e, phi)
		if gcd.Cmp(one) == 0 {
			return p, q, n, phi
		}
	}
}

func Generate() (*big.Int, *big.Int, *big.Int) {
	e := big.NewInt(65537)
	_, _, n, phi := generatePQ()

	if n == nil || phi == nil {
		return nil, nil, nil
	}

	d := new(big.Int).ModInverse(e, phi)
	if d == nil {
		return nil, nil, nil
	}

	return n, e, d
}
