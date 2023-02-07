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
	"runtime"
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
	optName = [...]string{"sortU4/F4", "lsw-U4/F4", "sortS", "lsw-S"}

	optFn = [...]func() float64{
		// optimize for arithmetic types
		func() float64 { return sumDurU4(false) + sumDurF4(false) },

		// optimize for lesswap sort
		func() float64 { return sumDurLswU4(false) + sumDurLswF4(false) },

		// optimize for string
		func() float64 { return sumDurS(false) },

		// optimize for lesswap sort (string key)
		func() float64 { return sumDurLswS(false) }}

	optInd int
)

func optStep(x, y int) float64 {
	if x <= 2*nsShort || y <= 2*x {
		return 9e9 // keep parameters feasible
	}
	MaxLenIns, MaxLenInsFC, MaxLenRec, MaxLenRecFC = x, x, y, y
	return optFn[optInd]()
}

func optRun(suffix string, ins0, rec0 int) {
	fmt.Printf("\n%s\nMaxLenIns%s MaxLenRec%s:\n", optName[optInd], suffix, suffix)

	_, _, _, n := opt.FindMinTri(2, ins0, rec0, ins0/4, rec0/4, optStep, optPrint)
	fmt.Println(n, "calls")
}

// Optimize max slice lengths for insertion sort/recursion
// Takes a long time, run with -tags tuneparam
func TestOptimize(t *testing.T) {
	tsPtr = t

	maxMaxGor = uint64(runtime.NumCPU())
	if maxMaxGor <= 1 {
		t.Fatal("need multiple cores to tune")
	}
	if maxMaxGor > 4 {
		maxMaxGor = 4
	}
	fmt.Println("max MaxGor:", maxMaxGor)

	optRun("", 80, 600)

	for optInd = 1; optInd < len(optFn); optInd++ {
		optRun("FC", 40, 300)
	}
}
