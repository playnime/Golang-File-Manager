package myrsa

import (
	"crypto/rand"
	"math/big"
)

// MillerRabin проверяет простоту x с n раундами
// n - количество раундов (чем больше, тем выше точность)
// при n=20 вероятность ошибки < 4^(-20)
func MillerRabin(n int, x *big.Int) bool {
	// граничные случаи
	two := big.NewInt(2)
	three := big.NewInt(3)

	if x.Cmp(two) == 0 || x.Cmp(three) == 0 {
		return true
	}
	// чётные и меньше 2 — составные
	if x.Cmp(two) < 0 || new(big.Int).Mod(x, two).Sign() == 0 {
		return false
	}

	// шаг 1: представляем x-1 = 2^s * d
	one := big.NewInt(1)
	xm1 := new(big.Int).Sub(x, one) // x - 1

	s := 0
	d := new(big.Int).Set(xm1)
	for new(big.Int).Mod(d, two).Sign() == 0 { // пока d чётное
		d.Rsh(d, 1) // d = d / 2
		s++
	}

	// шаг 2: n раундов с разными основаниями a
	for i := 0; i < n; i++ {
		// случайное a в диапазоне [2, x-2]
		xm2 := new(big.Int).Sub(x, two)
		a, err := rand.Int(rand.Reader, xm2)
		if err != nil {
			return false
		}
		a.Add(a, two) // a ∈ [2, x-2]

		// вычисляем a^d mod x
		y := new(big.Int).Exp(a, d, x)

		// если y == 1 или y == x-1 — этот раунд прошёл
		if y.Cmp(one) == 0 || y.Cmp(xm1) == 0 {
			continue
		}

		// возводим в квадрат s-1 раз
		composite := true
		for r := 0; r < s-1; r++ {
			y.Exp(y, two, x) // y = y^2 mod x
			if y.Cmp(xm1) == 0 {
				composite = false
				break
			}
		}

		if composite {
			return false // точно составное
		}
	}

	return true // вероятно простое
}
