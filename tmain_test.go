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
	"testing"
	"time"
	"unsafe"

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

func printSec(testName string, d time.Duration) float64 {
	sec := d.Seconds()
	fmt.Printf("%10s %5.2fs\n", testName, sec)
	return sec
}

// test & time sorting uint32 slices
// compare each result with standard sort.Slice
func TestUint(t *testing.T) {
	tsPtr = t

	medianCpstCompare("sort.Slice", nil, stdSort, false)
	sumDurU4(true) // sorty
	sumDurLswU4(true)
}

// test & time sorting float32 slices (NaNsmall)
// compare each result with standard sort.Slice
func TestFloatNaNsmall(t *testing.T) {
	tsPtr = t
	NaNoption = NaNsmall

	medianCpstCompare("sort.Slice", U4toF4, stdSort, false)
	sumDurF4(true) // sorty
	sumDurLswF4(true)
}

// test & time sorting float32 slices (NaNlarge)
// compare each result with standard sort.Slice
func TestFloatNaNlarge(t *testing.T) {
	tsPtr = t
	NaNoption = NaNlarge

	medianCpstCompare("sort.Slice", U4toF4, stdSort, false)
	sumDurF4(true) // sorty
	sumDurLswF4(true)
}

// test & time sorting string slices
// compare each result with standard sort.Slice
func TestString(t *testing.T) {
	tsPtr = t

	medianCpstCompare("sort.Slice", implantS, stdSort, false)
	sumDurS(true) // sorty
	sumDurLswS(true)
}

// test & time sorting []byte slices
// compare each result with standard sort.Slice
func TestByteSlice(t *testing.T) {
	tsPtr = t

	medianCpstCompare("sort.Slice", implantB, stdSort, false)
	sumDurB(true) // sorty
	sumDurLswB(true)
}

// test & time sorting string slices by length
// compare each result with standard sort.Slice
func TestStringByLen(t *testing.T) {
	tsPtr = t

	medianCpstCompare("sort.Slice", implantLenS, stdSortLen, false)
	sumDurLenS(true) // sorty
}

// test & time sorting []byte slices by length
// compare each result with standard sort.Slice
func TestByteSliceByLen(t *testing.T) {
	tsPtr = t

	medianCpstCompare("sort.Slice", implantLenB, stdSortLen, false)
	sumDurLenB(true) // sorty
}

func U4toU8(buf []uint32) interface{} {
	return sixb.U4toU8(buf)
}

func U4toI4(buf []uint32) interface{} {
	return *(*[]int32)(unsafe.Pointer(&buf))
}

func U4toI8(buf []uint32) interface{} {
	slc := sixb.U4toU8(buf)
	return *(*[]int64)(unsafe.Pointer(&slc))
}

func U4toF8(buf []uint32) interface{} {
	slc := sixb.U4toU8(buf)
	return *(*[]float64)(unsafe.Pointer(&slc))
}

func sortSignal(buf []uint32, prepare func([]uint32) interface{}, ch chan struct{}) {
	copyPrepSortTest(buf, prepare, SortSlice)
	if ch != nil {
		ch <- struct{}{}
	}
}

// is sorty multi-goroutine safe?
func TestConcurrent(t *testing.T) {
	tsPtr = t

	// buf1 & buf3 will get same random data in copyPrepSortTest, similarly buf2 & buf4
	buf1, buf2 := aaBuf[:bufHalf], aaBuf[bufHalf:]
	buf3, buf4 := bbBuf[:bufHalf], bbBuf[bufHalf:]

	lsPrep := [4]func([]uint32) interface{}{U4toU8, U4toI4, U4toI8, U4toF8}
	tch := make(chan struct{})

	for MaxGor = 2; MaxGor <= 4; MaxGor++ {
		for i := 0; i < len(lsPrep); i += 2 {

			fillSrc()
			go sortSignal(buf1, lsPrep[i], tch)
			go sortSignal(buf3, lsPrep[i], tch)
			go sortSignal(buf2, lsPrep[i+1], tch)
			sortSignal(buf4, lsPrep[i+1], nil)

			for k := 3; k > 0; k-- {
				<-tch // wait goroutines
			}
			compare(lsPrep[i](buf1), lsPrep[i](buf3))
			compare(lsPrep[i+1](buf2), lsPrep[i+1](buf4))
		}
	}
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
	if iArr[k-1] != 4 || iArr[k] != 5 || l != n {
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
