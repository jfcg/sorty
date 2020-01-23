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

// fill sort test for uint32
func fstUint(sd int64, ar []uint32, srt func([]uint32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if !IsSortedU4(ar) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

// fill sort test for uint64
func fstUint2(sd int64, ar []uint64, srt func([]uint64)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint64()
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if !IsSortedU8(ar) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

// fill sort test for int32
func fstInt(sd int64, ar []int32, srt func([]int32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = int32(rn.Uint32())
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if !IsSortedI4(ar) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

// fill sort test for int64
func fstInt2(sd int64, ar []int64, srt func([]int64)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = int64(rn.Uint64())
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if !IsSortedI8(ar) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

// fill sort test for float32
func fstFlt(sd int64, ar []float32, srt func([]float32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = float32(rn.NormFloat64())
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if !IsSortedF4(ar) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

// fill sort test for float64
func fstFlt2(sd int64, ar []float64, srt func([]float64)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.NormFloat64()
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if !IsSortedF8(ar) {
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
			ss[i].Data = unsafe.Pointer(&ar[k])
			ss[i].Len = 4
		}
	}
	return *(*[]string)(unsafe.Pointer(&ss)), ar[t:]
}

// fill sort test for string
func fstStr(sd int64, ar []uint32, srt func([]string)) time.Duration {
	as, ar := implant(ar, true)

	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(as)
	dur := time.Since(now)

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

// median of four durations
func medur(a, b, c, d time.Duration) time.Duration {
	if d < a {
		d, a = a, d
	}
	if b < a {
		b, a = a, b
	} else if d < b {
		d, b = b, d
	}
	if c < a {
		c, a = a, c
	} else if d < c {
		d, c = c, d
	}
	return (b + c) / 2
}

// median fst & compare for int
func mfcUint(tn string, srt func([]uint32), ar, ap []uint32) float64 {
	name = tn
	d1 := fstUint(1, ar, srt) // median of four different sorts
	d2 := fstUint(2, ar, srt)
	d3 := fstUint(3, ar, srt)
	d1 = medur(fstUint(4, ar, srt), d1, d2, d3)
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

func u2f(p *[]uint64) []float64 {
	return *(*[]float64)(unsafe.Pointer(p))
}

// median fst & compare for float
func mfcFlt(tn string, srt func([]float32), ar, ap []float32) float64 {
	name = tn
	d1 := fstFlt(5, ar, srt) // median of four different sorts
	d2 := fstFlt(6, ar, srt)
	d3 := fstFlt(7, ar, srt)
	d1 = medur(fstFlt(8, ar, srt), d1, d2, d3)
	compare(f2u(&ar), f2u(&ap))

	sec := d1.Seconds()
	if testing.Short() {
		fmt.Printf("%10s %5.2fs\n", name, sec)
	}
	return sec
}

// median fst & compare for string
func mfcStr(tn string, srt func([]string), ar, ap []uint32) float64 {
	name = tn
	d1 := fstStr(9, ar, srt) // median of four different sorts
	d2 := fstStr(10, ar, srt)
	d3 := fstStr(11, ar, srt)
	d1 = medur(fstStr(12, ar, srt), d1, d2, d3)

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

// return sum of SortU4() times for 2..4 goroutines
// compare with ap and among themselves
func sumtUint(ar, ap []uint32) float64 {
	s := .0
	for Mxg = 2; Mxg < 5; Mxg++ {
		srnm[6] = byte(Mxg + '0')
		s += mfcUint(string(srnm), SortU4, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of SortF4() times for 2..4 goroutines
// compare with ap and among themselves
func sumtFlt(ar, ap []float32) float64 {
	s := .0
	for Mxg = 2; Mxg < 5; Mxg++ {
		srnm[6] = byte(Mxg + '0')
		s += mfcFlt(string(srnm), SortF4, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of SortS() times for 2..4 goroutines
// compare with ap and among themselves
func sumtStr(ar, ap []uint32) float64 {
	s := .0
	for Mxg = 2; Mxg < 5; Mxg++ {
		srnm[6] = byte(Mxg + '0')
		s += mfcStr(string(srnm), SortS, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// int: return Sort(Col) time for 3 goroutines, compare with ap
func sumtCi(ar, ap []uint32) float64 {
	Mxg = 3 // sort via Collection
	return mfcUint("sorty-Col", func(aq []uint32) { Sort(uicol(aq)) }, ar, ap)
}

// int: return Sort2(Col2) time for 3 goroutines, compare with ap
func sumtC2i(ar, ap []uint32) float64 {
	Mxg = 3 // sort via Collection2
	return mfcUint("sorty-Col2", func(aq []uint32) { Sort2(uicol(aq)) }, ar, ap)
}

// float: return Sort(Col) time for 3 goroutines, compare with ap
func sumtCf(ar, ap []float32) float64 {
	Mxg = 3 // sort via Collection
	return mfcFlt("sorty-Col", func(aq []float32) { Sort(flcol(aq)) }, ar, ap)
}

// float: return Sort2(Col2) time for 3 goroutines, compare with ap
func sumtC2f(ar, ap []float32) float64 {
	Mxg = 3 // sort via Collection2
	return mfcFlt("sorty-Col2", func(aq []float32) { Sort2(flcol(aq)) }, ar, ap)
}

// string: return Sort2(Col2) time for 3 goroutines, compare with ap
func sumtC2s(ar, ap []uint32) float64 {
	Mxg = 3 // sort via Collection2
	return mfcStr("sorty-Col2", func(aq []string) { Sort2(stcol(aq)) }, ar, ap)
}

// sort int array with Sort3()
func sort3i(aq []uint32) {
	lsw := func(i, k, r, s int) bool {
		if aq[i] < aq[k] {
			if r != s {
				aq[r], aq[s] = aq[s], aq[r]
			}
			return true
		}
		return false
	}
	Sort3(len(aq), lsw)
}

// return Sort3() time for 3 goroutines, compare with ap
func sumtLi(ar, ap []uint32) float64 {
	Mxg = 3
	return mfcUint("sorty-lsw", sort3i, ar, ap)
}

// sort float array with Sort3()
func sort3f(aq []float32) {
	lsw := func(i, k, r, s int) bool {
		if aq[i] < aq[k] {
			if r != s {
				aq[r], aq[s] = aq[s], aq[r]
			}
			return true
		}
		return false
	}
	Sort3(len(aq), lsw)
}

// return Sort3() time for 3 goroutines, compare with ap
func sumtLf(ar, ap []float32) float64 {
	Mxg = 3
	return mfcFlt("sorty-lsw", sort3f, ar, ap)
}

// sort string array with Sort3()
func sort3s(aq []string) {
	lsw := func(i, k, r, s int) bool {
		if aq[i] < aq[k] {
			if r != s {
				aq[r], aq[s] = aq[s], aq[r]
			}
			return true
		}
		return false
	}
	Sort3(len(aq), lsw)
}

// return Sort3() time for 3 goroutines, compare with ap
func sumtLs(ar, ap []uint32) float64 {
	Mxg = 3
	return mfcStr("sorty-lsw", sort3s, ar, ap)
}

// types satisfying Collection* interfaces
type uicol []uint32
type flcol []float32
type stcol []string

func (c uicol) Len() int           { return len(c) }
func (c uicol) Less(i, k int) bool { return c[i] < c[k] }
func (c uicol) Swap(i, k int)      { c[i], c[k] = c[k], c[i] }
func (c uicol) LessSwap(i, k, r, s int) bool {
	if c[i] < c[k] {
		c[r], c[s] = c[s], c[r]
		return true
	}
	return false
}

func (c flcol) Len() int           { return len(c) }
func (c flcol) Less(i, k int) bool { return c[i] < c[k] }
func (c flcol) Swap(i, k int)      { c[i], c[k] = c[k], c[i] }
func (c flcol) LessSwap(i, k, r, s int) bool {
	if c[i] < c[k] {
		c[r], c[s] = c[s], c[r]
		return true
	}
	return false
}

func (c stcol) Len() int           { return len(c) }
func (c stcol) Less(i, k int) bool { return c[i] < c[k] }
func (c stcol) Swap(i, k int)      { c[i], c[k] = c[k], c[i] }
func (c stcol) LessSwap(i, k, r, s int) bool {
	if c[i] < c[k] {
		c[r], c[s] = c[s], c[r]
		return true
	}
	return false
}

// main test routine, needs -short flag
func TestShort(t *testing.T) {
	if !testing.Short() {
		t.SkipNow()
	}
	tst = t

	as := make([]float32, N)
	aq := make([]float32, N)
	ar, ap := f2u(&as), f2u(&aq)

	fmt.Println("Sorting uint32")
	mfcUint("sort.Slice", func(al []uint32) {
		sort.Slice(al, func(i, k int) bool { return al[i] < al[k] })
	}, ar, nil)
	mfcUint("sortutil", sortutil.Uint32s, ap, ar)
	mfcUint("zermelo", zuint32.Sort, ap, ar)
	sumtUint(ap, ar) // sorty
	sumtCi(ap, ar)
	sumtC2i(ap, ar)
	sumtLi(ap, ar)

	if !IsSorted(uicol(ap)) {
		t.Fatal("IsSorted() does not work")
	}
	if !IsSorted3(len(ap), func(i, k int) bool { return ap[i] < ap[k] }) {
		t.Fatal("IsSorted3() does not work")
	}

	fmt.Println("\nSorting float32")
	mfcFlt("sort.Slice", func(al []float32) {
		sort.Slice(al, func(i, k int) bool { return al[i] < al[k] })
	}, aq, nil)
	mfcFlt("sortutil", sortutil.Float32s, as, aq)
	mfcFlt("zermelo", zfloat32.Sort, as, aq)
	sumtFlt(as, aq) // sorty
	sumtCf(as, aq)
	sumtC2f(as, aq)
	sumtLf(as, aq)

	if !IsSorted(flcol(as)) {
		t.Fatal("IsSorted() does not work")
	}
	if !IsSorted3(len(as), func(i, k int) bool { return as[i] < as[k] }) {
		t.Fatal("IsSorted3() does not work")
	}

	fmt.Println("\nSorting string")
	mfcStr("sort.Slice", func(al []string) {
		sort.Slice(al, func(i, k int) bool { return al[i] < al[k] })
	}, ar, nil)
	mfcStr("sortutil", sortutil.Strings, ap, ar)
	mfcStr("radix", radix.Sort, ap, ar)
	sumtStr(ap, ar) // sorty
	sumtC2s(ap, ar)
	sumtLs(ap, ar)

	// Is Sort*() multi-goroutine safe?
	fmt.Println("\nConcurrent calls to Sort*()")
	name = "multi"
	K, ch := N/4, make(chan bool, 1)
	Mxg = 2
	aj, ak := sixb.I4tI8(ar), sixb.I4tI8(ap)
	am := u2f(&ak) // aj,am: uint64 & float64 buffers

	sasu := func(sd int64, al []uint64) {
		fstUint2(sd, al, SortU8) // sort and signal
		ch <- false
	}
	sasf := func(sd int64, al []float64) {
		fstFlt2(sd, al, SortF8)
		ch <- false
	}
	go sasu(21, aj[:K])
	go sasf(22, am[:K])
	go sasu(21, aj[K:])
	fstFlt2(22, am[K:], SortF8)

	for i := 3; i > 0; i-- {
		<-ch // wait others
	}
	K *= 2
	compare(ar[:K], ar[K:]) // same buffers
	compare(ap[:K], ap[K:])

	ao := *(*[]int64)(unsafe.Pointer(&aj)) // int64/32 buffers
	an := *(*[]int32)(unsafe.Pointer(&ap))

	sasi := func(sd int64, al []int32) {
		fstInt(sd, al, SortI4) // sort and signal
		ch <- false
	}
	sasj := func(sd int64, al []int64) {
		fstInt2(sd, al, SortI8)
		ch <- false
	}
	go sasi(23, an[:K])
	go sasj(24, ao[:N/4])
	go sasi(23, an[K:])
	fstInt2(24, ao[N/4:], SortI8)

	for i := 3; i > 0; i-- {
		<-ch // wait others
	}
	compare(ar[:K], ar[K:]) // same buffers
	compare(ap[:K], ap[K:])

	// SortI() calls SortI4() (on 32-bit) or SortI8() (on 64-bit).
	name = "SortI"
	SortI(iar)
	if !IsSortedI(iar) {
		t.Fatal("SortI() does not work")
	}

	// test Search()
	name = "Search"
	n := len(iar)
	k := Search(n, func(i int) bool { return iar[i] >= 5 })
	l := Search(n, func(i int) bool { return iar[i] >= 10 })
	if k <= 0 || k >= n || iar[k] != 5 || iar[k-1] != 4 || l != n {
		t.Fatal("Search() does not work")
	}
}

var iar = []int{
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2, 0, -1, 1, 2, 0,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2, 0, -1,
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 0, -1,
	-9, 8, -7, 6, -5, 4, -3, 2, -1, 0, 9, -8, 7, -6, 5, -4, 3, -2, 1, 0, 1, 2, 0, -1,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, -9}

// Optimize max array lengths for insertion sort/recursion (Mli,Mlr)
// takes a long time, run without -short flag
func TestOpt(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tst = t

	pro := func(x, y int, v float64) { // print optimum
		fmt.Printf("%3d %3d %5.2fs\n", x, y, v)
	}
	as := make([]float32, N)
	aq := make([]float32, 0, N)
	ar, ap := f2u(&as), f2u(&aq)

	name := []string{"U4/F4", "S", "", "2", "3", "3s"}
	fn := []func() float64{
		// optimize for native arithmetic types
		func() float64 { return sumtUint(ar, ap) + sumtFlt(as, aq) },

		// optimize for native string
		func() float64 { return sumtStr(ar, ap) },

		// optimize for Collection interface
		func() float64 { return sumtCi(ar, ap) + sumtCf(as, aq) },

		// optimize for Collection2 interface
		func() float64 { return sumtC2i(ar, ap) + sumtC2f(as, aq) },

		// optimize for function-based interface
		func() float64 { return sumtLi(ar, ap) + sumtLf(as, aq) },

		// optimize for function-based interface (string key)
		func() float64 { return sumtLs(ar, ap) }}

	for i := 0; i < len(fn); i++ {
		fmt.Println("\nSort" + name[i])

		_, _, _, n := opt.FindMinTri(2, 96, 449, 12, 64,
			func(x, y int) float64 {
				Mli, Mlr = x, y
				return fn[i]()
			}, pro)
		fmt.Println(n, "calls")
	}
}
