//go:build !tuneparam
// +build !tuneparam

/*	Copyright (c) 2021, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/jfcg/sixb"
)

func TestMinMax(t *testing.T) {
	for n := uint(2); n <= 3*nsConc; n++ {
		for slen := 2 * n; slen <= 3*MaxLenRec; slen++ {

			first, step, last := minMaxSample(slen, n)
			diff := last - first

			// basic checks
			if !(first < last && last < slen && 2 <= step && step <= diff) {
				t.Fatal("must have first < last < slen, 2 ≤ step ≤ last-first")
			}

			// first is larger tail, tails differ by at most 1
			if !(slen-last-1 <= first && first <= slen-last) {
				t.Fatal("must have slen-last-1 ≤ first ≤ slen-last")
			}

			// equidistant
			if (n-1)*step != diff {
				t.Fatal("must have first + (n-1) * step = last")
			}

			tail := slen - diff // 1 + #members in both tails

			// max distance to non-selected members optimal?
			if tail >= n && first > (step+1)>>1 || // from tails to inside
				step>>1 > (tail+n-1)>>1 { // from inside to tails
				t.Fatal("max distance sub-optimal")
			}
		}
	}

	for slen := uint(8); slen <= 3*MaxLenRec; slen++ {
		f1, s1, _ := minMaxSample(slen, 4)
		f2, s2 := minMaxFour(uint32(slen))

		if f1 != uint(f2) || s1 != uint(s2) {
			t.Fatal("minMaxFour != minMaxSample")
		}
	}
}

func printSec(tn string, d time.Duration) float64 {
	sec := d.Seconds()
	fmt.Printf("%10s %5.2fs\n", tn, sec)
	return sec
}

// sort and signal
func sasU8(sd int64, al []uint64, ch chan struct{}) {
	fstU8(sd, al, sortU8)
	ch <- struct{}{}
}

func sasF8(sd int64, al []float64, ch chan struct{}) {
	fstF8(sd, al, sortF8)
	ch <- struct{}{}
}

func sasI4(sd int64, al []int32, ch chan struct{}) {
	fstI4(sd, al, sortI4)
	ch <- struct{}{}
}

func sasI8(sd int64, al []int64, ch chan struct{}) {
	fstI8(sd, al, sortI8)
	ch <- struct{}{}
}

// test & time sorting uint32 slices, compare their results
func TestUint(t *testing.T) {
	tsPtr = t

	mfcU4("sort.Slice", func(al []uint32) {
		sort.Slice(al, func(i, k int) bool { return al[i] < al[k] })
	}, bufbu, nil)
	sumtU4(bufau, bufbu) // sorty
	sumtLswU4(bufau, bufbu)

	if IsSorted(len(bufau), func(i, k, r, s int) bool { return bufau[i] < bufau[k] }) != 0 {
		t.Fatal("IsSorted() does not work")
	}
}

// test & time sorting float32 slices, compare their results
func TestFloat(t *testing.T) {
	tsPtr = t

	mfcF4("sort.Slice", func(al []float32) {
		sort.Slice(al, func(i, k int) bool { return al[i] < al[k] })
	}, bufbf, nil)
	sumtF4(bufaf, bufbf) // sorty
	sumtLswF4(bufaf, bufbf)

	if IsSorted(len(bufaf), func(i, k, r, s int) bool { return bufaf[i] < bufaf[k] }) != 0 {
		t.Fatal("IsSorted() does not work")
	}
}

// test & time sorting string slices, compare their results
func TestString(t *testing.T) {
	tsPtr = t

	mfcS("sort.Slice", func(al []string) {
		sort.Slice(al, func(i, k int) bool { return al[i] < al[k] })
	}, bufbu, nil)
	sumtS(bufau, bufbu) // sorty
	sumtLswS(bufau, bufbu)
}

// test & time sorting []byte slices, compare their results
func TestByteSlice(t *testing.T) {
	tsPtr = t

	mfcB("sort.Slice", func(al [][]byte) {
		sort.Slice(al, func(i, k int) bool { return sixb.BtoS(al[i]) < sixb.BtoS(al[k]) })
	}, bufbu, nil)
	sumtB(bufau, bufbu) // sorty
}

// test & time sorting string slices 'by length', compare their results
func TestStringByLen(t *testing.T) {
	tsPtr = t

	mfcLenS("sort.Slice", func(al []string) {
		sort.Slice(al, func(i, k int) bool { return len(al[i]) < len(al[k]) })
	}, bufbu, nil)
	sumtLenS(bufau, bufbu) // sorty
}

// test & time sorting []byte slices 'by length', compare their results
func TestByteSliceByLen(t *testing.T) {
	tsPtr = t

	mfcLenB("sort.Slice", func(al [][]byte) {
		sort.Slice(al, func(i, k int) bool { return len(al[i]) < len(al[k]) })
	}, bufbu, nil)
	sumtLenB(bufau, bufbu) // sorty
}

// Is Sort*() multi-goroutine safe?
func TestConcurrent(t *testing.T) {
	tsPtr = t

	bufK, bufL, tch := bufN/2, bufN/4, make(chan struct{})
	MaxGor = 2

	// two concurrent calls to sortU8() & sortF8() each
	// up to 8 goroutines total
	go sasU8(96, bufbu2[:bufL:bufL], tch)
	go sasF8(97, bufaf2[:bufL:bufL], tch)
	go sasU8(96, bufbu2[bufL:], tch)
	fstF8(97, bufaf2[bufL:], sortF8)

	for i := 3; i > 0; i-- {
		<-tch // wait others
	}
	compareU4(bufbu[:bufK:bufK], bufbu[bufK:]) // same buffers
	compareU4(bufau[:bufK:bufK], bufau[bufK:])

	// two concurrent calls to sortI4() & sortI8() each
	// up to 8 goroutines total
	go sasI4(98, bufai[:bufK:bufK], tch)
	go sasI8(99, bufbi2[:bufL:bufL], tch)
	go sasI4(98, bufai[bufK:], tch)
	fstI8(99, bufbi2[bufL:], sortI8)

	for i := 3; i > 0; i-- {
		<-tch // wait others
	}
	compareU4(bufbu[:bufK:bufK], bufbu[bufK:]) // same buffers
	compareU4(bufau[:bufK:bufK], bufau[bufK:])
}

// Sort()ing short slices
func TestShort(t *testing.T) {
	tsPtr = t

	for l := -3; l < 2; l++ {
		Sort(l, iarlsw)
		if iArr[0] != 9 || iArr[1] != 8 {
			t.Fatal("Sort()ing short slices does not work")
		}
	}
	for l := 2; l < 4; l++ {
		Sort(l, iarlsw)
		for k := 2; k >= 0; k-- {
			if iArr[k] != iArr[12+k-l] {
				t.Fatal("Sort()ing short slices does not work")
			}
		}
	}

	// SortSlice() calls sortI4() on 32-bit, sortI8() on 64-bit
	SortSlice(iArr)
	if IsSortedSlice(iArr) != 0 {
		t.Fatal("SortSlice/IsSortedSlice does not work")
	}

	// test Search()
	n := len(iArr)
	k := Search(n, func(i int) bool { return iArr[i] >= 5 })
	l := Search(n, func(i int) bool { return iArr[i] >= 10 })
	if k <= 0 || k >= n || iArr[k] != 5 || iArr[k-1] != 4 || l != n {
		t.Fatal("Search() does not work")
	}
}

var iArr = []int{
	9, 8, 7, 6, 5, 4, 3, 2, 1, 7, 8, 9, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2, 0, -1, 1, 2, 0,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2, 0, -1,
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 0, -1,
	-9, 8, -7, 6, -5, 4, -3, 2, -1, 0, 9, -8, 7, -6, 5, -4, 3, -2, 1, 0, 1, 2, 0, -1,
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 0, -1,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, -9}

func iarlsw(i, k, r, s int) bool {
	if iArr[i] < iArr[k] {
		if r != s {
			iArr[r], iArr[s] = iArr[s], iArr[r]
		}
		return true
	}
	return false
}
