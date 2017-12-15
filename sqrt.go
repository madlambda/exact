package exact

// Sqrt computes the square root of x.
// The precision is DefPrecision.
// To use arbitrary precision, see Sqrtp.
// Note: this function is orders of magnitude slow than math.Sqrt,
// but the result can be made as precise as you wish by using
// Sqrtp.
func Sqrt(x Frac) Frac {
	return Sqrtp(x, DefPrecision)
}

// Sqrtp computes the square root of x using precision prec.
// Note: Depending on the value of prec, this function can
// take hours, days, years, a life time, ..., to complete.
func Sqrtp(x, prec Frac) Frac {
	if x.P.Cmp(zero) == 0 {
		return Zero()
	}

	g := One() // first guess

	for !closeEnough(Div(x, g), g, prec) {
		g = betterGuess(x, g)
	}

	return g
}

// make g = (g+(x/g))/2
func betterGuess(x, g Frac) Frac {
	return Mul(Add(g, Div(x, g)), NewFrac(1, 2, false))
}

// abs(a-b) < precision
func closeEnough(a, b, prec Frac) bool {
	return Lt(Abs(Sub(a, b)), prec)
}
