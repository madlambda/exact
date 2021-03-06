package exact_test

import (
	"math"
	"testing"

	"github.com/madlambda/exact"
)

func almostFloat(x, y, ε float64) bool {
	return math.Abs(x-y) <= ε
}

func assert(t *testing.T, b bool, msg string) {
	t.Helper()
	if !b {
		t.Fatal(msg)
	}
}

func assertAlmost(t *testing.T, x, y, ε float64, msg string) {
	t.Helper()
	assert(t, almostFloat(x, y, ε), fmt("Fail: %s. Differs: %.12f != %.12f",
		msg, x, y))
}

func TestRatSqrtAgainstFloat(t *testing.T) {
	for i := uint64(0); i < 1000; i++ {
		fp := float64(i)
		frac := exact.NewRat(i, 1)

		sfp := math.Sqrt(fp)
		sf := exact.Sqrt(frac)

		assertAlmost(t, sfp, sf.Inexact(),
			exact.DefPrecision.Inexact(), fmt("sqrt(%d)", i))
	}
}

func BenchmarkExactSqrt2(b *testing.B) {
	two := exact.NewRat(2, 1)
	for n := 0; n < b.N; n++ {
		exact.Sqrt(two)
	}
}

func BenchmarkMathSqrt2(b *testing.B) {
	// unfair comparison, but nice to see the difference
	for n := 0; n < b.N; n++ {
		math.Sqrt(2.0)
	}
}
