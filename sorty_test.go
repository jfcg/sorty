/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"fmt"
	"github.com/jfcg/opt"
	"github.com/jfcg/sixb"
	"github.com/shawnsmithdev/zermelo/zfloat32"
	"github.com/shawnsmithdev/zermelo/zuint32"
	"github.com/twotwotwo/sorts/sortutil"
	"github.com/yourbasic/radix"
	"math/rand"
	"sort"
	"testing"
	"time"
	"unsafe"
)

const N = 1 << 26

var tst *testing.T
var name string

// fill sort test
func fst(sd int64, ar []uint32, srt func([]uint32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(ar)
	dur := time.Now().Sub(now)

	if !IsSortedU4(ar) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

// fill sort test
func fst2(sd int64, ar []float32, srt func([]float32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = float32(rn.NormFloat64())
	}

	now := time.Now()
	srt(ar)
	dur := time.Now().Sub(now)

	if !IsSortedF4(ar) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

// implant strings into ar
func implant(ar []uint32, fill bool) ([]string, []uint32) {
	// string size is 4*t bytes
	t := int(unsafe.Sizeof("") >> 2)

	// ar will hold n strings (headers followed by 4-byte bodies)
	n := len(ar) / (t + 1)

	t *= n // total string headers space
	ss := sixb.I4tSs(ar[:t])

	if fill {
		for i, k := n-1, len(ar)-1; i >= 0; i, k = i-1, k-1 {
			ss[i].Data = uintptr(unsafe.Pointer(&ar[k]))
			ss[i].Len = 4
		}
	}
	return *(*[]string)(unsafe.Pointer(&ss)), ar[t:]
}

// fill sort test for string
func fst3(sd int64, ar []uint32, srt func([]string)) time.Duration {
	as, ar := implant(ar, true)

	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(as)
	dur := time.Now().Sub(now)

	if !IsSortedS(as) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

func compare(ar, ap []uint32) {
	l := len(ap)
	if l <= 0 {
		return
	}
	if len(ar) != l {
		tst.Fatal(name, "length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if ar[i] != ap[i] {
			tst.Fatal(name, "values mismatch:", i, ar[i], ap[i])
		}
	}
}

func compareS(ar, ap []string) {
	l := len(ap)
	if len(ar) != l {
		tst.Fatal(name, "length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if ar[i] != ap[i] {
			tst.Fatal(name, "values mismatch:", i, ar[i], ap[i])
		}
	}
}

// median of three
func medur(a, b, c time.Duration) time.Duration {
	if a > b {
		a, b = b, a
	}
	if c <= a {
		return a
	}
	if c >= b {
		return b
	}
	return c
}

// median fst & compare
func mfc(srt func([]uint32), ar, ap []uint32) float64 {
	d1 := fst(1, ar, srt) // median of three different sorts
	d2 := fst(2, ar, srt)
	d1 = medur(fst(3, ar, srt), d1, d2)
	compare(ar, ap)

	sec := d1.Seconds()
	if testing.Short() {
		fmt.Printf("%10s %5.2fs\n", name, sec)
	}
	return sec
}

func f2u(p *[]float32) []uint32 {
	return *(*[]uint32)(unsafe.Pointer(p))
}

// median fst & compare
func mfc2(srt func([]float32), ar, ap []float32) float64 {
	d1 := fst2(4, ar, srt) // median of three different sorts
	d2 := fst2(5, ar, srt)
	d1 = medur(fst2(6, ar, srt), d1, d2)
	compare(f2u(&ar), f2u(&ap))

	sec := d1.Seconds()
	if testing.Short() {
		fmt.Printf("%10s %5.2fs\n", name, sec)
	}
	return sec
}

// median fst & compare for string
func mfc3(srt func([]string), ar, ap []uint32) float64 {
	d1 := fst3(7, ar, srt) // median of three different sorts
	d2 := fst3(8, ar, srt)
	d1 = medur(fst3(9, ar, srt), d1, d2)

	if len(ap) > 0 {
		as, ar := implant(ar, false)
		aq, ap := implant(ap, false)
		compareS(as, aq)
		compare(ar, ap)
	}

	sec := d1.Seconds()
	if testing.Short() {
		fmt.Printf("%10s %5.2fs\n", name, sec)
	}
	return sec
}

var srnm = []byte("sorty-0")

// return sum of Sort*() times for 2..4 goroutines & Sort(Col)
// compare with ap and among themselves
func sumt(ar, ap []uint32) float64 {
	s := .0
	for Mxg = 2; Mxg < 5; Mxg++ {
		srnm[6] = byte(Mxg + '0')
		name = string(srnm)
		s += mfc(SortU4, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}

	Mxg = 3
	name = "Sort(Col)" // sort via Collection
	s += mfc(func(aq []uint32) { Sort(uicol(aq)) }, ar, ap)
	return s
}

// return sum of Sort*() times for 2..4 goroutines & Sort(Col)
// compare with ap and among themselves
func sumt2(ar, ap []float32) float64 {
	s := .0
	for Mxg = 2; Mxg < 5; Mxg++ {
		srnm[6] = byte(Mxg + '0')
		name = string(srnm)
		s += mfc2(SortF4, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}

	Mxg = 3
	name = "Sort(Col)" // sort via Collection
	s += mfc2(func(aq []float32) { Sort(flcol(aq)) }, ar, ap)
	return s
}

// return sum of SortS() times for 2..4 goroutines
// compare with ap and among themselves
func sumt3(ar, ap []uint32) float64 {
	s := .0
	for Mxg = 2; Mxg < 5; Mxg++ {
		srnm[6] = byte(Mxg + '0')
		name = string(srnm)
		s += mfc3(SortS, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

type uicol []uint32
type flcol []float32

func (c uicol) Len() int           { return len(c) }
func (c uicol) Less(i, k int) bool { return c[i] < c[k] }
func (c uicol) Swap(i, k int)      { c[i], c[k] = c[k], c[i] }

func (c flcol) Len() int           { return len(c) }
func (c flcol) Less(i, k int) bool { return c[i] < c[k] }
func (c flcol) Swap(i, k int)      { c[i], c[k] = c[k], c[i] }

func TestShort(t *testing.T) {
	if !testing.Short() {
		t.SkipNow()
	}
	tst = t

	as := make([]float32, N)
	aq := make([]float32, N)
	ar, ap := f2u(&as), f2u(&aq)

	fmt.Println("Sorting uint32")
	name = "sort.Slice"
	mfc(func(ar []uint32) {
		sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] })
	}, ar, nil)
	name = "sortutil"
	mfc(sortutil.Uint32s, ap, ar)
	name = "zermelo"
	mfc(zuint32.Sort, ap, ar)
	sumt(ap, ar) // sorty
	if !IsSorted(uicol(ap)) {
		t.Fatal("IsSorted() does not work")
	}

	fmt.Println("\nSorting float32")
	name = "sort.Slice"
	mfc2(func(ar []float32) {
		sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] })
	}, aq, nil)
	name = "sortutil"
	mfc2(sortutil.Float32s, as, aq)
	name = "zermelo"
	mfc2(zfloat32.Sort, as, aq)
	sumt2(as, aq) // sorty
	if !IsSorted(flcol(as)) {
		t.Fatal("IsSorted() does not work")
	}

	fmt.Println("\nSorting string")
	name = "sort.Slice"
	mfc3(func(ar []string) {
		sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] })
	}, ar, nil)
	name = "sortutil"
	mfc3(sortutil.Strings, ap, ar)
	name = "radix"
	mfc3(radix.Sort, ap, ar)
	sumt3(ap, ar) // sorty

	// Is Sort*() multi-goroutine safe?
	fmt.Println("\nConcurrent calls to SortU4()")
	name = "multi"
	K, ch := N/2, make(chan bool, 1)
	sas := func(sd int64, ar []uint32) {
		fst(sd, ar, SortU4) // SortU4 and signal
		ch <- false
	}

	Mxg = 2
	go sas(19, ar[:K])
	go sas(20, ar[K:])
	go sas(19, ap[:K])
	fst(20, ap[K:], SortU4)

	for i := 3; i > 0; i-- {
		<-ch // wait others
	}
	compare(ar[:K], ap[:K])
	compare(ar[K:], ap[K:])

	// SortI calls SortI4 (on 32-bit) or SortI8 (on 64-bit).
	name = "SortI"
	SortI(iar)
	if !IsSortedI(iar) {
		t.Fatal("SortI does not work")
	}
}

var iar = []int{
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2, 0, -1, 1, 2, 0,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2, 0, -1,
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 0, -1,
	-9, 8, -7, 6, -5, 4, -3, 2, -1, 0, 9, -8, 7, -6, 5, -4, 3, -2, 1, 0, 1, 2, 0, -1,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, -9}

// Optimize max array lengths for insertion sort/recursion (Mli,Mlr)
func TestOpt(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tst = t

	pro := func(x, y int, v float64) { // print optimum
		fmt.Printf("%d %d %.2fs\n", x, y, v)
	}
	as := make([]float32, N)
	aq := make([]float32, 0, N)
	ar, ap := f2u(&as), f2u(&aq)

	_, _, _, n := opt.FindMinTri(2, 80, 901, 16, 128, func(x, y int) float64 {
		Mli, Mlr = x, y
		return sumt(ar, ap) + sumt2(as, aq) + sumt3(ar, ap)
	}, pro)
	fmt.Println("made", n, "calls")
}
