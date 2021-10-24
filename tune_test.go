//go:build tuneparam
// +build tuneparam

/*	Copyright (c) 2021, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"fmt"
	"testing"
	"time"

	"github.com/jfcg/opt"
)

func printSec(_ string, d time.Duration) float64 {
	return d.Seconds()
}

func printOpt(x, y int, v float64) {
	fmt.Printf("%3d %3d %5.2fs\n", x, y, v)
}

var (
	optName = [...]string{"SortU4/F4", "Lsw-U4/F4", "SortS", "Lsw-S"}

	bufap = make([]uint32, N)

	optFn = [...]func() float64{
		// optimize for native arithmetic types
		func() float64 { return sumtU4(bufau, bufap[:0]) + sumtF4(bufaf, bufbf[:0]) },

		// optimize for function-based sort
		// carry over bufap,bufbf for further comparison
		func() float64 { return sumtLswU4(bufau, bufap) + sumtLswF4(bufaf, bufbf) },

		// optimize for native string
		func() float64 { return sumtS(bufau, bufap[:0]) },

		// optimize for function-based sort (string key)
		// carry over bufap for further comparison
		func() float64 { return sumtLswS(bufau, bufap) }}
)

// Optimize max slice lengths for insertion sort/recursion
// Takes a long time, run with -tags tuneparam
func TestOptimize(t *testing.T) {
	tsPtr = t

	s1, s2 := "MaxLenIns", 96

	for i := 0; i < len(optFn); i++ {
		fmt.Printf("\n%s\n%s MaxLenRec:\n", optName[i], s1)

		_, _, _, n := opt.FindMinTri(2, s2, 480, 16, 120,
			func(x, y int) float64 {
				if x < 10 || y <= 2*x {
					return 9e9 // keep parameters feasible
				}
				MaxLenIns, MaxLenInsFC, MaxLenRec = x, x, y
				return optFn[i]()
			}, printOpt)
		fmt.Println(n, "calls")

		s1, s2 = "MaxLenInsFC", 48
	}
}
