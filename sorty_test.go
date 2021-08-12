/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/jfcg/sixb"
)

const N = 1 << 26

var tsPtr *testing.T
var tsName string

// fill sort test for uint32
func fstU4(sd int64, ar []uint32, srt func([]uint32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if IsSortedU4(ar) != 0 {
		tsPtr.Fatal(tsName, "not sorted")
	}
	return dur
}

// fill sort test for uint64
func fstU8(sd int64, ar []uint64, srt func([]uint64)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint64()
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if IsSortedU8(ar) != 0 {
		tsPtr.Fatal(tsName, "not sorted")
	}
	return dur
}

// fill sort test for int32
func fstI4(sd int64, ar []int32, srt func([]int32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = int32(rn.Uint32())
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if IsSortedI4(ar) != 0 {
		tsPtr.Fatal(tsName, "not sorted")
	}
	return dur
}

// fill sort test for int64
func fstI8(sd int64, ar []int64, srt func([]int64)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = int64(rn.Uint64())
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if IsSortedI8(ar) != 0 {
		tsPtr.Fatal(tsName, "not sorted")
	}
	return dur
}

// fill sort test for float32
func fstF4(sd int64, ar []float32, srt func([]float32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = float32(rn.NormFloat64())
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if IsSortedF4(ar) != 0 {
		tsPtr.Fatal(tsName, "not sorted")
	}
	return dur
}

// fill sort test for float64
func fstF8(sd int64, ar []float64, srt func([]float64)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.NormFloat64()
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if IsSortedF8(ar) != 0 {
		tsPtr.Fatal(tsName, "not sorted")
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
	ss := sixb.I4tSs(ar[:t:t])

	if fill {
		for i, k := n-1, len(ar)-1; i >= 0; i, k = i-1, k-1 {
			ss[i].Data = unsafe.Pointer(&ar[k])
			ss[i].Len = 4
		}
	}
	return *(*[]string)(unsafe.Pointer(&ss)), ar[t:]
}

// fill sort test for string
func fstS(sd int64, ar []uint32, srt func([]string)) time.Duration {
	as, ar := implant(ar, true)

	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(as)
	dur := time.Since(now)

	if IsSortedS(as) != 0 {
		tsPtr.Fatal(tsName, "not sorted")
	}
	return dur
}

func compareU4(ar, ap []uint32) {
	l := len(ap)
	if l <= 0 {
		return
	}
	if len(ar) != l {
		tsPtr.Fatal(tsName, "length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if ar[i] != ap[i] {
			tsPtr.Fatal(tsName, "values mismatch:", i, ar[i], ap[i])
		}
	}
}

func compareS(ar, ap []string) {
	l := len(ap)
	if len(ar) != l {
		tsPtr.Fatal(tsName, "length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if ar[i] != ap[i] {
			tsPtr.Fatal(tsName, "values mismatch:", i, ar[i], ap[i])
		}
	}
}

// median of four durations
func medur(a, b, c, d time.Duration) time.Duration {
	if d < b {
		d, b = b, d
	}
	if c < a {
		c, a = a, c
	}
	if d < c {
		c = d
	}
	if b < a {
		b = a
	}
	return (b + c) >> 1
}

// median fst & compare for uint32
func mfcU4(tn string, srt func([]uint32), ar, ap []uint32) float64 {
	tsName = tn
	d1 := fstU4(1, ar, srt) // median of four different sorts
	d2 := fstU4(2, ar, srt)
	d3 := fstU4(3, ar, srt)
	d1 = medur(fstU4(4, ar, srt), d1, d2, d3)

	compareU4(ar, ap)
	return printSec(d1)
}

// slice conversions
func F4toU4(p *[]float32) []uint32 {
	return *(*[]uint32)(unsafe.Pointer(p))
}

func F4toI4(p *[]float32) []int32 {
	return *(*[]int32)(unsafe.Pointer(p))
}

func U8toF8(p *[]uint64) []float64 {
	return *(*[]float64)(unsafe.Pointer(p))
}

func U8toI8(p *[]uint64) []int64 {
	return *(*[]int64)(unsafe.Pointer(p))
}

// median fst & compare for float32
func mfcF4(tn string, srt func([]float32), ar, ap []float32) float64 {
	tsName = tn
	d1 := fstF4(5, ar, srt) // median of four different sorts
	d2 := fstF4(6, ar, srt)
	d3 := fstF4(7, ar, srt)
	d1 = medur(fstF4(8, ar, srt), d1, d2, d3)

	compareU4(F4toU4(&ar), F4toU4(&ap))
	return printSec(d1)
}

// median fst & compare for string
func mfcS(tn string, srt func([]string), ar, ap []uint32) float64 {
	tsName = tn
	d1 := fstS(9, ar, srt) // median of four different sorts
	d2 := fstS(10, ar, srt)
	d3 := fstS(11, ar, srt)
	d1 = medur(fstS(12, ar, srt), d1, d2, d3)

	if len(ap) > 0 {
		as, ar := implant(ar, false)
		aq, ap := implant(ap, false)
		compareS(as, aq)
		compareU4(ar, ap)
	}
	return printSec(d1)
}

var srtName = []byte("sorty-0")

// return sum of SortU4() times for 1..4 goroutines
// compare with ap and among themselves
func sumtU4(ar, ap []uint32) float64 {
	s := .0
	for Mxg = 1; Mxg < 5; Mxg++ {
		srtName[6] = byte(Mxg + '0')
		s += mfcU4(string(srtName), SortU4, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of SortF4() times for 1..4 goroutines
// compare with ap and among themselves
func sumtF4(ar, ap []float32) float64 {
	s := .0
	for Mxg = 1; Mxg < 5; Mxg++ {
		srtName[6] = byte(Mxg + '0')
		s += mfcF4(string(srtName), SortF4, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of SortS() times for 1..4 goroutines
// compare with ap and among themselves
func sumtS(ar, ap []uint32) float64 {
	s := .0
	for Mxg = 1; Mxg < 5; Mxg++ {
		srtName[6] = byte(Mxg + '0')
		s += mfcS(string(srtName), SortS, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// sort uint32 array with Sort()
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
	Sort(len(aq), lsw)
}

var lswName = []byte("sortyLsw-0")

// return sum of sort3i() times for 1..4 goroutines
// compare with ap and among themselves
func sumtLswU4(ar, ap []uint32) float64 {
	s := .0
	for Mxg = 1; Mxg < 5; Mxg++ {
		lswName[9] = byte(Mxg + '0')
		s += mfcU4(string(lswName), sort3i, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// sort float32 array with Sort()
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
	Sort(len(aq), lsw)
}

// return sum of sort3f() times for 1..4 goroutines
// compare with ap and among themselves
func sumtLswF4(ar, ap []float32) float64 {
	s := .0
	for Mxg = 1; Mxg < 5; Mxg++ {
		lswName[9] = byte(Mxg + '0')
		s += mfcF4(string(lswName), sort3f, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// sort string array with Sort()
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
	Sort(len(aq), lsw)
}

// return sum of sort3s() times for 1..4 goroutines
// compare with ap and among themselves
func sumtLswS(ar, ap []uint32) float64 {
	s := .0
	for Mxg = 1; Mxg < 5; Mxg++ {
		lswName[9] = byte(Mxg + '0')
		s += mfcS(string(lswName), sort3s, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}
