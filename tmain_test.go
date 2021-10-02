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
	//"github.com/shawnsmithdev/zermelo/zfloat32"
	//"github.com/shawnsmithdev/zermelo/zuint32"
	//"github.com/twotwotwo/sorts/sortutil"
	//"github.com/yourbasic/radix"
)

func printSec(tn string, d time.Duration) float64 {
	sec := d.Seconds()
	fmt.Printf("%10s %5.2fs\n", tn, sec)
	return sec
}

// sort and signal
func sasU8(sd int64, al []uint64, ch chan bool) {
	fstU8(sd, al, SortU8)
	ch <- false
}

func sasF8(sd int64, al []float64, ch chan bool) {
	fstF8(sd, al, SortF8)
	ch <- false
}

func sasI4(sd int64, al []int32, ch chan bool) {
	fstI4(sd, al, SortI4)
	ch <- false
}

func sasI8(sd int64, al []int64, ch chan bool) {
	fstI8(sd, al, SortI8)
	ch <- false
}

// test & time sorting uint32 slices for different libraries, compare their results
func TestUint(t *testing.T) {
	tsPtr = t

	fmt.Println("\nSorting []uint32")
	mfcU4("sort.Slice", func(al []uint32) {
		sort.Slice(al, func(i, k int) bool { return al[i] < al[k] })
	}, bufbu, nil)
	//mfcU4("sortutil", sortutil.Uint32s, au, bu)
	//mfcU4("zermelo", zuint32.Sort, au, bu)
	sumtU4(bufau, bufbu) // sorty
	sumtLswU4(bufau, bufbu)

	if IsSorted(len(bufau), func(i, k, r, s int) bool { return bufau[i] < bufau[k] }) != 0 {
		t.Fatal("IsSorted() does not work")
	}
}

// test & time sorting float32 slices for different libraries, compare their results
func TestFloat(t *testing.T) {
	tsPtr = t

	fmt.Println("\nSorting []float32")
	mfcF4("sort.Slice", func(al []float32) {
		sort.Slice(al, func(i, k int) bool { return al[i] < al[k] })
	}, bufbf, nil)
	//mfcF4("sortutil", sortutil.Float32s, af, bf)
	//mfcF4("zermelo", zfloat32.Sort, af, bf)
	sumtF4(bufaf, bufbf) // sorty
	sumtLswF4(bufaf, bufbf)

	if IsSorted(len(bufaf), func(i, k, r, s int) bool { return bufaf[i] < bufaf[k] }) != 0 {
		t.Fatal("IsSorted() does not work")
	}
}

// test & time sorting string/[]byte slices for
// different libraries, compare their results
func TestString(t *testing.T) {
	tsPtr = t

	fmt.Println("\nSorting []string")
	mfcS("sort.Slice", func(al []string) {
		sort.Slice(al, func(i, k int) bool { return al[i] < al[k] })
	}, bufbu, nil)
	//mfcS("sortutil", sortutil.Strings, au, bu)
	//mfcS("radix", radix.Sort, au, bu)
	sumtS(bufau, bufbu) // sorty
	sumtLswS(bufau, bufbu)

	fmt.Println("\nSorting [][]byte")
	mfcB("sort.Slice", func(al [][]byte) {
		sort.Slice(al, func(i, k int) bool { return sixb.BtoS(al[i]) < sixb.BtoS(al[k]) })
	}, bufbu, nil)
	sumtB(bufau, bufbu) // sorty
}

// test & time sorting string/[]byte slices 'by length'
// for different libraries, compare their results
func TestLength(t *testing.T) {
	tsPtr = t

	fmt.Println("\nSorting []string by length")
	mfcLenS("sort.Slice", func(al []string) {
		sort.Slice(al, func(i, k int) bool { return len(al[i]) < len(al[k]) })
	}, bufbu, nil)
	sumtLenS(bufau, bufbu) // sorty

	fmt.Println("\nSorting [][]byte by length")
	mfcLenB("sort.Slice", func(al [][]byte) {
		sort.Slice(al, func(i, k int) bool { return len(al[i]) < len(al[k]) })
	}, bufbu, nil)
	sumtLenB(bufau, bufbu) // sorty
}

// Is Sort*() multi-goroutine safe?
func TestConcurrent(t *testing.T) {
	tsPtr = t

	fmt.Println("\nConcurrent calls to Sort*()")
	K, L, ch := N/2, N/4, make(chan bool)
	Mxg = 2

	// two concurrent calls to SortU8() & SortF8() each
	// up to 8 goroutines total
	go sasU8(96, bufbu2[:L:L], ch)
	go sasF8(97, bufaf2[:L:L], ch)
	go sasU8(96, bufbu2[L:], ch)
	fstF8(97, bufaf2[L:], SortF8)

	for i := 3; i > 0; i-- {
		<-ch // wait others
	}
	compareU4(bufbu[:K:K], bufbu[K:]) // same buffers
	compareU4(bufau[:K:K], bufau[K:])

	// two concurrent calls to SortI4() & SortI8() each
	// up to 8 goroutines total
	go sasI4(98, bufai[:K:K], ch)
	go sasI8(99, bufbi2[:L:L], ch)
	go sasI4(98, bufai[K:], ch)
	fstI8(99, bufbi2[L:], SortI8)

	for i := 3; i > 0; i-- {
		<-ch // wait others
	}
	compareU4(bufbu[:K:K], bufbu[K:]) // same buffers
	compareU4(bufau[:K:K], bufau[K:])
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

	// SortI() calls SortI4() (on 32-bit) or SortI8() (on 64-bit).
	SortI(iArr)
	if IsSortedI(iArr) != 0 {
		t.Fatal("SortI() does not work")
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
