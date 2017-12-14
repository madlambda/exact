package exact_test

import (
	"fmt"
	"testing"

	"github.com/madlambda/exact"
)

type Frac = exact.Frac // just because I'm lazy

var testcases = []struct {
	a, b, sum, sub, mul, div Frac
}{
	{
		a:   Frac{P: 1, Q: 1},
		b:   Frac{P: 1, Q: 1},
		sum: Frac{P: 2, Q: 1},
		sub: Frac{P: 0, Q: 1},
		div: Frac{P: 1, Q: 1},
		mul: Frac{P: 1, Q: 1},
	},
	{
		a:   Frac{P: 1, Q: 2},
		b:   Frac{P: 1, Q: 2},
		sum: Frac{P: 1, Q: 1},
		sub: Frac{P: 0, Q: 1},
		div: Frac{P: 1, Q: 1},
		mul: Frac{P: 1, Q: 4},
	},
	{
		a:   Frac{P: 10, Q: 2},
		b:   Frac{P: 15, Q: 3},
		sum: Frac{P: 10, Q: 1},
		sub: Frac{P: 0, Q: 1},
		div: Frac{P: 1, Q: 1},
		mul: Frac{P: 25, Q: 1},
	},
}

func assertEqual(t *testing.T, a, b Frac) {
	if !exact.Cmp(a, b) {
		t.Fatalf("fail: %s != %s", a, b)
	}
}

func TestAdd(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("%s+%s", tc.a, tc.b), func(t *testing.T) {
			r := exact.Add(tc.a, tc.b)
			assertEqual(t, r, tc.sum)
		})
	}
}

func TestSub(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("%s-%s", tc.a, tc.b), func(t *testing.T) {
			r := exact.Sub(tc.a, tc.b)
			assertEqual(t, r, tc.sub)
		})
	}
}

func TestDiv(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("%s/%s", tc.a, tc.b), func(t *testing.T) {
			r := exact.Div(tc.a, tc.b)
			assertEqual(t, r, tc.div)
		})
	}
}

func TestMul(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("%s*%s", tc.a, tc.b), func(t *testing.T) {
			r := exact.Mul(tc.a, tc.b)
			assertEqual(t, r, tc.mul)
		})
	}
}
