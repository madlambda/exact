// Package exact implements common math for rational numbers.
// You must use this package if you need to avoid floating point
// roudoff errors.
package exact

import (
	"math/big"
)

type (
	// Rat is a fraction
	Rat struct {
		Sign bool
		P    *big.Int // P is the numerator
		Q    *big.Int // Q is the denominator
	}
)

const (
	positive = false
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)

	// DefPrecision is the default precision used when
	// the function does not specify one.
	DefPrecision Rat
)

func init() {
	var prec, _, _ = big.NewFloat(0).Parse("1.0e100", 10)
	z := big.NewInt(0)
	prec.Int(z)
	DefPrecision = NewBigRat(big.NewInt(1), z)
}

// One is the whole.
func One() Rat { return NewBigRat(one, one) }

// Zero is the empty.
func Zero() Rat { return NewBigRat(zero, one) }

// NewRat creates a new positive rational number using p, q as numerator and
// denominator, respectively.
func NewRat(p, q uint64) Rat {
	return Rat{
		Sign: positive,
		P:    big.NewInt(0).SetUint64(p),
		Q:    big.NewInt(0).SetUint64(q),
	}
}

// NewNegRat creates a new negative rational number using p, q as numerator and
// denominator, respectively.
func NewNegRat(p, q uint64) Rat {
	return Rat{
		Sign: !positive,
		P:    big.NewInt(0).SetUint64(p),
		Q:    big.NewInt(0).SetUint64(q),
	}
}

// Inverse of the rational number
func (r Rat) Inverse() Rat {
	return Rat{
		Sign: r.Sign,
		P:    r.Q,
		Q:    r.P,
	}
}

// Neg returns the negative version of the number
func (r Rat) Neg() Rat {
	return NewNegBigRat(r.P, r.Q)
}

// NewBigRat creates a new positive rational number in the same way
// as NewRat but using big.Int as numerator and denominator.
func NewBigRat(p, q *big.Int) Rat {
	return newBigRat(p, q, positive)
}

// NewNegBigRat creates a new negative rational number in the same way
// as NewRat but using big.Int as numerator and denominator.
func NewNegBigRat(p, q *big.Int) Rat {
	return newBigRat(p, q, !positive)
}

func newBigRat(p, q *big.Int, sign bool) Rat {
	if q.Cmp(zero) == 0 {
		panic("division by zero")
	}
	return Rat{
		Sign: sign,
		P:    p,
		Q:    q,
	}
}

// IsZero tells if f is zero
func (r Rat) IsZero() bool {
	return r.P.Cmp(zero) == 0
}

// add but ignores sign
func add(a, b Rat) Rat {
	p1 := big.NewInt(0).Mul(a.P, b.Q)
	p2 := big.NewInt(0).Mul(b.P, a.Q)
	return Rat{
		P: big.NewInt(0).Add(p1, p2),
		Q: big.NewInt(0).Mul(a.Q, b.Q),
	}
}

// Add two rationale numbers a and b.
func Add(a, b Rat) Rat {
	if a.Sign == b.Sign {
		// -a + (-b) == (-a)-b == -(a+b)
		// +a + (+b) == +(a+b)
		r := add(NewBigRat(a.P, a.Q), NewBigRat(b.P, b.Q))
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
	return Rat{
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
func Sub(a, b Rat) Rat {
	// a - b = a + (-b)
	b.Sign = !b.Sign
	return Add(a, b)
}

// Mul multiplies a and b.
func Mul(a, b Rat) Rat {
	r := Rat{
		Sign: a.Sign && b.Sign,
		P:    big.NewInt(0).Mul(a.P, b.P),
		Q:    big.NewInt(0).Mul(a.Q, b.Q),
	}

	return r
}

// Div divides a and b
func Div(a, b Rat) Rat {
	if b.P.Cmp(zero) == 0 || b.Q.Cmp(zero) == 0 {
		panic("division by zero")
	}
	// a = p/q
	// b = p'/q'
	// a/b = (p/q)/(p'/q') = (p/q)*(q'/p')
	return Mul(a, Rat{
		Sign: a.Sign && b.Sign,
		P:    b.Q,
		Q:    b.P,
	})
}

// Abs returns the absolute value of x
func Abs(x Rat) Rat {
	if x.Sign {
		x.Sign = false
		return x
	}
	return x
}

// Lt is the less than (<) comparator.
func Lt(a, b Rat) bool {
	a, b = a.Simplify(), b.Simplify()
	if a.Sign == b.Sign &&
		a.Q.Cmp(b.Q) == 0 {
		return a.P.Cmp(b.P) < 0
	}

	r := Sub(a, b)
	return r.Sign
}

// Cmp compares if a equals b.
func Cmp(a, b Rat) bool {
	if a.P.Cmp(zero) == 0 &&
		b.P.Cmp(zero) == 0 {
		return true
	}
	as := a.Simplify()
	bs := b.Simplify()
	return as.P.Cmp(bs.P) == 0 &&
		as.Q.Cmp(bs.Q) == 0
}

// String returns the string representation of the rational number.
func (r Rat) String() string {
	slash := rune('/')

	var digits []rune
	p := r.P.Text(10)
	if r.P.Sign() < 0 {
		digits = append(digits, rune('-'))
	}

	digits = append(digits, []rune(p)...)
	digits = append(digits, slash)
	digits = append(digits, []rune(r.Q.Text(10))...)
	return string(digits)
}

// Inexact returns the inexact floating point from the fraction.
func (r Rat) Inexact() float64 {
	v := new(big.Rat)
	v.SetFrac(r.P, r.Q)
	fp, _ := v.Float64()
	return fp
}

// Simplify fraction
func (r Rat) Simplify() Rat {
	if r.P.Cmp(zero) == 0 || r.P.Cmp(one) == 0 {
		return r
	}
	cd := big.NewInt(0).Abs(r.P).GCD(nil, nil, r.P, r.Q) // common divisor
	return Rat{
		P: big.NewInt(0).Div(r.P, cd),
		Q: big.NewInt(0).Div(r.Q, cd),
	}
}
