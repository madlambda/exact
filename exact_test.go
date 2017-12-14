package exact_test

import (
	format "fmt"
	"math/big"
	"testing"

	"github.com/madlambda/exact"
)

func newInt(t string) *big.Int {
	v, _ := big.NewInt(0).SetString(t, 10)
	return v
}

var (
	sqrt11_3 = exact.NewFrac2(
		newInt(
			"8907430012405601974227841812502266124751907164389458124457923971110416952138563811933935901308162252289",
		),
		newInt(
			"4651753611446999121325528753135336086170729864116047599840048338398572005441776052283614001901065410989",
		),
		false,
	)

	sqrt10_3 = exact.NewFrac2(
		newInt(
			"182140598762941461355069160591888526271635523676933716216619188883432692646610141830040815824221130179601964717075465433643881983192710343327452700070368556752347372334044277761",
		),
		newInt(
			"99762514579960592434515661886422227213820543479328130216353111486818126195968731365445978030596906270110314538555708556984066735218679109730815587464861598957360538800363906816",
		),
		false,
	)
)

var fmt = format.Sprintf
var testcases = []struct {
	a, b, sum, sub, mul, div exact.Frac
	sqrtA, sqrtB             exact.Frac
	lt, gt                   bool
	abs                      exact.Frac
}{
	{
		a:     exact.One(),
		b:     exact.One(),
		sum:   exact.NewFrac(2, 1, false),
		sub:   exact.NewFrac(0, 1, false),
		div:   exact.One(),
		mul:   exact.One(),
		sqrtA: exact.One(),
		sqrtB: exact.One(),
		lt:    false,
	},
	{
		a: exact.NewFrac(1, 2, false),
		b: exact.NewFrac(1, 2, false),
		sqrtA: exact.NewFrac2(
			newInt(
				"48926646634423881954586808839856694558492182258668537145547700898547222910968507268117381704646657",
			),
			newInt(
				"69192727231838199530637090778029723034779720143976685296374209532493131389050939536650584353662464",
			),
			false,
		),
		sqrtB: exact.NewFrac2(
			newInt(
				"48926646634423881954586808839856694558492182258668537145547700898547222910968507268117381704646657",
			),
			newInt(
				"69192727231838199530637090778029723034779720143976685296374209532493131389050939536650584353662464",
			),
			false,
		),
		sum: exact.NewFrac(1, 1, false),
		sub: exact.NewFrac(0, 1, false),
		div: exact.NewFrac(1, 1, false),
		mul: exact.NewFrac(1, 4, false),
		lt:  false,
	},
	{
		a: exact.NewFrac(10, 2, false),
		b: exact.NewFrac(15, 3, false),
		sqrtA: exact.NewFrac2(
			newInt(
				"316837008400094222150776738483768236006420971486980607",
			),
			newInt(
				"141693817714056513234709965875411919657707794958199867",
			),
			false,
		),
		sqrtB: exact.NewFrac2(
			newInt(
				"316837008400094222150776738483768236006420971486980607",
			),
			newInt(
				"141693817714056513234709965875411919657707794958199867",
			),
			false,
		),
		sum: exact.NewFrac(10, 1, false),
		sub: exact.NewFrac(0, 1, false),
		div: exact.NewFrac(1, 1, false),
		mul: exact.NewFrac(25, 1, false),
		lt:  false,
	},
	{
		a:     exact.NewFrac(11, 3, false),
		b:     exact.NewFrac(10, 3, false),
		sqrtA: sqrt11_3,
		sqrtB: sqrt10_3,
		sum:   exact.NewFrac(21, 3, false),
		sub:   exact.NewFrac(1, 3, false),
		div:   exact.NewFrac(11, 10, false),
		mul:   exact.NewFrac(110, 9, false),
		lt:    false,
	},
	{
		a:     exact.NewFrac(10, 3, false),
		b:     exact.NewFrac(11, 3, false),
		sqrtA: sqrt10_3,
		sqrtB: sqrt11_3,
		sum:   exact.NewFrac(21, 3, false),
		sub:   exact.NewFrac(1, 3, true),
		div:   exact.NewFrac(10, 11, false),
		mul:   exact.NewFrac(110, 9, false),
		lt:    true,
	},
	{
		a:     exact.NewFrac(10, 3, true),
		b:     exact.NewFrac(11, 3, false),
		sqrtA: sqrt10_3,
		sqrtB: sqrt11_3,
		sum:   exact.NewFrac(1, 3, false),
		sub:   exact.NewFrac(21, 3, true),
		div:   exact.NewFrac(10, 11, true),
		mul:   exact.NewFrac(110, 9, true),
		lt:    true,
	},
}

func assertEqual(t *testing.T, a, b exact.Frac) {
	if !exact.Cmp(a, b) {
		t.Fatalf("fail: %s != %s", a, b)
	}
}

func TestAdd(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt("%s+%s", tc.a, tc.b), func(t *testing.T) {
			r := exact.Add(tc.a, tc.b)
			assertEqual(t, r, tc.sum)
		})
	}
}

func TestSub(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt("%s-%s", tc.a, tc.b), func(t *testing.T) {
			r := exact.Sub(tc.a, tc.b)
			assertEqual(t, r, tc.sub)
		})
	}
}

func TestDiv(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt("%s/%s", tc.a, tc.b), func(t *testing.T) {
			r := exact.Div(tc.a, tc.b)
			assertEqual(t, r, tc.div)
		})
	}
}

func TestMul(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt("%s*%s", tc.a, tc.b), func(t *testing.T) {
			r := exact.Mul(tc.a, tc.b)
			assertEqual(t, r, tc.mul)
		})
	}
}

func TestLt(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt("%s < %s", tc.a, tc.b), func(t *testing.T) {
			r := exact.Lt(tc.a, tc.b)
			if r != tc.lt {
				t.Fatalf("%s < %s != %v (got %v)", tc.a, tc.b, tc.lt, r)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt("sqrt(%s),sqrt(%s)", tc.a, tc.b), func(t *testing.T) {
			ra := exact.Sqrt(tc.a)
			rb := exact.Sqrt(tc.b)

			if !exact.Cmp(ra, tc.sqrtA) {
				raS := ra.Simplify()
				t.Fatalf("sqrt(%s) != %v, got %s (%f)", tc.a, tc.sqrtA, raS, raS.Inexact())
			}

			if !exact.Cmp(rb, tc.sqrtB) {
				rbS := rb.Simplify()
				t.Fatalf("sqrt(%s) != %v, got %s (%f)", tc.b, tc.sqrtB,
					rbS, rbS.Inexact())
			}
		})
	}
}
