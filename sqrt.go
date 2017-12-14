package exact

import (
	"math/big"
)

var SqrtPrec Frac

func init() {
	var prec, _, _ = big.NewFloat(0).Parse("1.0e100", 10)
	z := big.NewInt(0)
	prec.Int(z)
	SqrtPrec = NewFrac2(big.NewInt(1), z, false)
}

// Sqrt computes the square root of a.
// Note: slow, use with caution.
// TODO: make iterative
func Sqrt(x Frac) (ret Frac) {
	if x.P.Cmp(zero) == 0 {
		return Zero()
	}

	g := One() // first guess

	for !closeEnough(Div(x, g), g) {
		g = betterGuess(x, g)
	}

	return g
}

func betterGuess(x, g Frac) Frac {
	return Mul(Add(g, Div(x, g)), NewFrac(1, 2, false))
}

func closeEnough(a, b Frac) bool {
	// abs(a-b) < precision
	return Lt(Abs(Sub(a, b)), SqrtPrec)
}
