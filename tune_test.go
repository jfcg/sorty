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

func optPrint(x, y int, v float64) {
	fmt.Printf("%3d %3d %5.2fs\n", x, y, v)
}

var (
	optName = [...]string{"SortU4/F4", "Lsw-U4/F4", "SortS", "Lsw-S"}

	bufap = make([]uint32, bufN)

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

	optInd int
)

func optStep(x, y int) float64 {
	if x <= 2*nsShort || y <= 2*x {
		return 9e9 // keep parameters feasible
	}
	MaxLenIns, MaxLenInsFC, MaxLenRec = x, x, y
	return optFn[optInd]()
}

func optRun(name string, ins0, rec0 int) {
	fmt.Printf("\n%s\n%s MaxLenRec:\n", optName[optInd], name)

	_, _, _, n := opt.FindMinTri(2, ins0, rec0, ins0>>2, rec0>>2, optStep, optPrint)
	fmt.Println(n, "calls")
}

// Optimize max slice lengths for insertion sort/recursion
// Takes a long time, run with -tags tuneparam
func TestOptimize(t *testing.T) {
	tsPtr = t

	optRun("MaxLenIns", 96, 480)

	for optInd = 1; optInd < len(optFn); optInd++ {
		optRun("MaxLenInsFC", 48, 240)
	}
}
