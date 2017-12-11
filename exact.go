// Package exact implements common math for rational numbers.
// You must use this package if you need to avoid floating point
// roudoff errors.
package exact

type (
	// Frac is a fraction
	Frac struct {
		P int64 // P is the numerator
		Q int64 // Q is the denominator
	}
)

// Add two rationale numbers a and b.
func Add(a, b Frac) Frac {
	return Frac{
		P: a.P*b.Q + b.P*a.Q,
		Q: a.Q * b.Q,
	}
}

// Sub subtract the rational numbers a and b.
func Sub(a, b Frac) Frac {
	return Add(a, Frac{
		P: -b.P,
		Q: b.Q,
	})
}

// Mul multiplies a and b.
func Mul(a, b Frac) Frac {
	return Frac{
		P: a.P * b.P,
		Q: a.Q * b.Q,
	}
}

// Div divides a and b
func Div(a, b Frac) Frac {
	return Mul(a, Frac{
		P: b.Q,
		Q: b.P,
	})
}

// gcd is the greatest common dividor
func gcd(a, b int64) int64 {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

// Simplify the fraction
func Simplify(a Frac) Frac {
	cd := gcd(a.P, a.Q) // common divisor
	return Frac{
		P: a.P / cd,
		Q: a.Q / cd,
	}
}

// Cmp compares if a equals b.
func Cmp(a, b Frac) bool {
	as := Simplify(a)
	bs := Simplify(b)
	return as.P == bs.P &&
		as.Q == bs.Q
}

func reverse(s []rune) {
	if len(s) < 2 {
		return
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// integer to []rune
func itors(i int64, runeBase rune) []rune {
	if i == 0 {
		return []rune{runeBase}
	}

	digits := []rune{}
	for i > 0 {
		digits = append(digits, runeBase+(rune(i%10)))
		i = i / 10
	}
	reverse(digits)
	return digits
}

func (f Frac) String() string {
	slash := rune(0x2044)
	digits := itors(f.P, rune(0x2070))
	digits = append(digits, slash)
	digits = append(digits, itors(f.Q, 0x2080)...)
	return string(digits)
}

// Inexact returns the inexact floating point from a fraction.
func Inexact(a Frac) float64 {
	return float64(a.P) / float64(a.Q)
}
