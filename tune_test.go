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

func printSec(d time.Duration) float64 {
	return d.Seconds()
}

func printOpt(x, y int, v float64) {
	fmt.Printf("%3d %3d %5.2fs\n", x, y, v)
}

// Optimize max array lengths for insertion sort/recursion (Mli,Mlr)
// Takes a long time, run with -tags tuneparam
func TestOpt(t *testing.T) {
	tsPtr = t

	as := make([]float32, N)
	aq := make([]float32, N)
	ar := F4toU4(&as)
	ap := make([]uint32, N)

	nm := [...]string{"SortU4/F4", "Lsw-U4/F4", "SortS", "Lsw-S"}
	fn := [...]func() float64{
		// optimize for native arithmetic types
		func() float64 { return sumtU4(ar, ap[:0]) + sumtF4(as, aq[:0]) },

		// optimize for function-based sort
		// carry over ap,aq for further comparison
		func() float64 { return sumtLswU4(ar, ap) + sumtLswF4(as, aq) },

		// optimize for native string
		func() float64 { return sumtS(ar, ap[:0]) },

		// optimize for function-based sort (string key)
		// carry over ap for further comparison
		func() float64 { return sumtLswS(ar, ap) }}

	s1, s2 := "Mli", 96

	for i := 0; i < len(fn); i++ {
		fmt.Printf("\n%s\n%s Mlr:\n", nm[i], s1)

		_, _, _, n := opt.FindMinTri(2, s2, 480, 16, 120,
			func(x, y int) float64 {
				Mli, Hmli, Mlr = x, x, y
				return fn[i]()
			}, printOpt)
		fmt.Println(n, "calls")

		s1, s2 = "Hmli", 48
	}
}
