// Package exact implements common math for rational numbers.
// You must use this package if you need to avoid floating point
// roudoff errors.
package exact

import (
	"math/big"
)

type (
	// Frac is a fraction
	Frac struct {
		Sign bool
		P    *big.Int // P is the numerator
		Q    *big.Int // Q is the denominator
	}
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)

	// DefPrecision is the default precision used when
	// the function does not specify one.
	DefPrecision Frac
)

func init() {
	var prec, _, _ = big.NewFloat(0).Parse("1.0e100", 10)
	z := big.NewInt(0)
	prec.Int(z)
	DefPrecision = NewFrac2(big.NewInt(1), z, false)
}

// One is the whole.
func One() Frac { return NewFrac2(one, one, false) }

// Zero is the empty.
func Zero() Frac { return NewFrac2(zero, one, false) }

// NewFrac creates a new fraction using p and q as numerator and
// denominator, respectively. Sign tells if it's a negative fraction.
func NewFrac(p, q uint64, sign bool) Frac {
	return Frac{
		Sign: sign,
		P:    big.NewInt(0).SetUint64(p),
		Q:    big.NewInt(0).SetUint64(q),
	}
}

// NewFrac2 creates a new fraction in the same way as NewFrac
// but using big.Int as numerator and denominator.
func NewFrac2(p, q *big.Int, sign bool) Frac {
	if q.Cmp(zero) == 0 {
		panic("division by zero")
	}
	return Frac{
		Sign: sign,
		P:    p,
		Q:    q,
	}
}

// IsZero tells if f is zero
func (f Frac) IsZero() bool {
	return f.P.Cmp(zero) == 0
}

// add but ignores sign
func add(a, b Frac) Frac {
	p1 := big.NewInt(0).Mul(a.P, b.Q)
	p2 := big.NewInt(0).Mul(b.P, a.Q)
	return Frac{
		P: big.NewInt(0).Add(p1, p2),
		Q: big.NewInt(0).Mul(a.Q, b.Q),
	}
}

// Add two rationale numbers a and b.
func Add(a, b Frac) Frac {
	if a.Sign == b.Sign {
		// -a + (-b) == (-a)-b == -(a+b)
		// +a + (+b) == +(a+b)
		r := add(NewFrac2(a.P, a.Q, false), NewFrac2(b.P, b.Q, false))
		r.Sign = a.Sign
		return r
	}

	at := big.NewInt(0).Mul(a.P, b.Q)
	bt := big.NewInt(0).Mul(b.P, a.Q)

	p := big.NewInt(0)
	var sign bool

	m := max(at, bt)
	if m == at {
		p = p.Sub(at, bt)
		sign = a.Sign
	} else {
		p = p.Sub(bt, at)
		sign = b.Sign
	}
	return Frac{
		Sign: sign,
		P:    p,
		Q:    big.NewInt(0).Mul(a.Q, b.Q),
	}
}

func max(a, b *big.Int) *big.Int {
	if a.Cmp(b) > 0 {
		return a
	}
	return b
}

// Sub subtract the rational numbers a and b.
func Sub(a, b Frac) Frac {
	// a - b = a + (-b)
	b.Sign = !b.Sign
	return Add(a, b)
}

// Mul multiplies a and b.
func Mul(a, b Frac) Frac {
	r := Frac{
		Sign: a.Sign && b.Sign,
		P:    big.NewInt(0).Mul(a.P, b.P),
		Q:    big.NewInt(0).Mul(a.Q, b.Q),
	}

	return r
}

// Div divides a and b
func Div(a, b Frac) Frac {
	if b.P.Cmp(zero) == 0 || b.Q.Cmp(zero) == 0 {
		panic("division by zero")
	}
	// a = p/q
	// b = p'/q'
	// a/b = (p/q)/(p'/q') = (p/q)*(q'/p')
	return Mul(a, Frac{
		Sign: a.Sign && b.Sign,
		P:    b.Q,
		Q:    b.P,
	})
}

// Abs returns the absolute value of x
func Abs(x Frac) Frac {
	if x.Sign {
		x.Sign = false
		return x
	}
	return x
}

// Lt is the less than (<) comparator.
func Lt(a, b Frac) bool {
	a, b = a.Simplify(), b.Simplify()
	if a.Sign == b.Sign &&
		a.Q.Cmp(b.Q) == 0 {
		return a.P.Cmp(b.P) < 0
	}

	r := Sub(a, b)
	return r.Sign
}

// Cmp compares if a equals b.
func Cmp(a, b Frac) bool {
	if a.P.Cmp(zero) == 0 &&
		b.P.Cmp(zero) == 0 {
		return true
	}
	as := a.Simplify()
	bs := b.Simplify()
	return as.P.Cmp(bs.P) == 0 &&
		as.Q.Cmp(bs.Q) == 0
}

// text to []rune
func textToRunes(num string) []rune {
	digits := []rune{}
	for _, s := range num {
		digits = append(digits, s)
	}

	return digits
}

func (f Frac) String() string {
	slash := rune('/')

	var digits []rune
	p := f.P.Text(10)
	if f.P.Sign() < 0 {
		digits = append(digits, rune('-'))
	}

	digits = append(digits, textToRunes(p)...)
	digits = append(digits, slash)
	digits = append(digits, textToRunes(f.Q.Text(10))...)
	return string(digits)
}

// Inexact returns the inexact floating point from the fraction.
func (f Frac) Inexact() float64 {
	v := new(big.Rat)
	v.SetFrac(f.P, f.Q)
	fp, _ := v.Float64()
	return fp
}

// Simplify fraction
func (f Frac) Simplify() Frac {
	if f.P.Cmp(zero) == 0 || f.P.Cmp(one) == 0 {
		return f
	}
	cd := big.NewInt(0).Abs(f.P).GCD(nil, nil, f.P, f.Q) // common divisor
	return Frac{
		P: big.NewInt(0).Div(f.P, cd),
		Q: big.NewInt(0).Div(f.Q, cd),
	}
}
