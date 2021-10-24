/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/jfcg/sixb"
)

const N = 1 << 26

var (
	// a & b buffers will hold all slices to sort
	bufaf = make([]float32, N)
	bufbf = make([]float32, N)

	// different type views of the same buffers
	bufau  = F4toU4(&bufaf) // uint32
	bufbu  = F4toU4(&bufbf)
	bufai  = F4toI4(&bufaf)     // int32
	bufau2 = sixb.U4toU8(bufau) // uint64
	bufbu2 = sixb.U4toU8(bufbu)
	bufaf2 = U8toF8(&bufau2) // float64
	bufbi2 = U8toI8(&bufbu2) // int64

	tsPtr *testing.T
)

// fill sort test for uint32
func fstU4(sd int64, ar []uint32, srt func([]uint32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(ar)
	dur := time.Since(now)

	if isSortedU4(ar) != 0 {
		tsPtr.Fatal("not sorted")
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

	if isSortedU8(ar) != 0 {
		tsPtr.Fatal("not sorted")
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

	if isSortedI4(ar) != 0 {
		tsPtr.Fatal("not sorted")
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

	if isSortedI8(ar) != 0 {
		tsPtr.Fatal("not sorted")
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

	if isSortedF4(ar) != 0 {
		tsPtr.Fatal("not sorted")
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

	if isSortedF8(ar) != 0 {
		tsPtr.Fatal("not sorted")
	}
	return dur
}

// implant strings into ar
func implantS(ar []uint32, fill bool) ([]string, []uint32) {
	// string size is 4*t bytes
	t := sixb.StrSize >> 2

	// ar will hold n strings (headers followed by 4-byte bodies)
	n := len(ar) / (t + 1)

	t *= n // total string headers space
	ss := sixb.U4toStrs(ar[:t:t])

	if fill {
		for k := len(ar); n > 0; {
			n--
			k--
			ss[n].Data = unsafe.Pointer(&ar[k])
			ss[n].Len = 4
		}
	}
	return *(*[]string)(unsafe.Pointer(&ss)), ar[t:]
}

// fill sort test for string
func fstS(sd int64, ar []uint32, srt func([]string)) time.Duration {
	as, ar := implantS(ar, true)

	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(as)
	dur := time.Since(now)

	if isSortedS(as) != 0 {
		tsPtr.Fatal("not sorted")
	}
	return dur
}

// implant []byte's into ar
func implantB(ar []uint32, fill bool) ([][]byte, []uint32) {
	// []byte size is 4*t bytes
	t := sixb.SliceSize >> 2

	// ar will hold n []byte's (headers followed by 4-byte bodies)
	n := len(ar) / (t + 1)

	t *= n // total []byte headers space
	bs := sixb.U4toSlcs(ar[:t:t])

	if fill {
		for k := len(ar); n > 0; {
			n--
			k--
			bs[n].Data = unsafe.Pointer(&ar[k])
			bs[n].Len = 4
			bs[n].Cap = 4
		}
	}
	return *(*[][]byte)(unsafe.Pointer(&bs)), ar[t:]
}

// fill sort test for []byte
func fstB(sd int64, ar []uint32, srt func([][]byte)) time.Duration {
	ab, ar := implantB(ar, true)

	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(ab)
	dur := time.Since(now)

	if isSortedB(ab) != 0 {
		tsPtr.Fatal("not sorted")
	}
	return dur
}

// implant strings into ar (SortLen)
func implantLenS(sd int64, ar []uint32, fill bool) []string {
	// string size is 4*t bytes
	t := sixb.StrSize >> 2

	// ar will hold n string headers
	n := len(ar) / t

	t *= n // total string headers space
	ss := sixb.U4toStrs(ar[:t:t])

	if fill {
		rn := rand.New(rand.NewSource(sd))

		for L := 4*len(ar) + 1; n > 0; {
			n--
			// string bodies start at &ar[0] with random lengths up to 4*len(ar) bytes
			ss[n].Data = unsafe.Pointer(&ar[0])
			ss[n].Len = rn.Intn(L)
		}
	}
	return *(*[]string)(unsafe.Pointer(&ss))
}

// fill sort test for string (SortLen)
func fstLenS(sd int64, ar []uint32, srt func([]string)) time.Duration {
	as := implantLenS(sd, ar, true)

	now := time.Now()
	srt(as)
	dur := time.Since(now)

	if IsSortedLen(as) != 0 {
		tsPtr.Fatal("not sorted")
	}
	return dur
}

// implant []byte's into ar (SortLen)
func implantLenB(sd int64, ar []uint32, fill bool) [][]byte {
	// []byte size is 4*t bytes
	t := sixb.SliceSize >> 2

	// ar will hold n []byte headers
	n := len(ar) / t

	t *= n // total []byte headers space
	bs := sixb.U4toSlcs(ar[:t:t])

	if fill {
		rn := rand.New(rand.NewSource(sd))

		for L := 4*len(ar) + 1; n > 0; {
			n--
			// []byte bodies start at &ar[0] with random lengths up to 4*len(ar) bytes
			bs[n].Data = unsafe.Pointer(&ar[0])
			bs[n].Len = rn.Intn(L)
			bs[n].Cap = bs[n].Len
		}
	}
	return *(*[][]byte)(unsafe.Pointer(&bs))
}

// fill sort test for []byte (SortLen)
func fstLenB(sd int64, ar []uint32, srt func([][]byte)) time.Duration {
	ab := implantLenB(sd, ar, true)

	now := time.Now()
	srt(ab)
	dur := time.Since(now)

	if IsSortedLen(ab) != 0 {
		tsPtr.Fatal("not sorted")
	}
	return dur
}

func compareU4(ar, ap []uint32) {
	l := len(ap)
	if l <= 0 {
		return
	}
	if len(ar) != l {
		tsPtr.Fatal("length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if ar[i] != ap[i] {
			tsPtr.Fatal("values mismatch:", i, ar[i], ap[i])
		}
	}
}

func compareS(ar, ap []string) {
	l := len(ap)
	if len(ar) != l {
		tsPtr.Fatal("length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if ar[i] != ap[i] {
			tsPtr.Fatal("values mismatch:", i, ar[i], ap[i])
		}
	}
}

func compareB(ar, ap [][]byte) {
	l := len(ap)
	if len(ar) != l {
		tsPtr.Fatal("length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if sixb.BtoS(ar[i]) != sixb.BtoS(ap[i]) {
			tsPtr.Fatal("values mismatch:", i, ar[i], ap[i])
		}
	}
}

func compareLenS(ar, ap []string) {
	l := len(ap)
	if len(ar) != l {
		tsPtr.Fatal("length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if len(ar[i]) != len(ap[i]) {
			tsPtr.Fatal("lengths mismatch:", i, len(ar[i]), len(ap[i]))
		}
	}
}

func compareLenB(ar, ap [][]byte) {
	l := len(ap)
	if len(ar) != l {
		tsPtr.Fatal("length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if len(ar[i]) != len(ap[i]) {
			tsPtr.Fatal("lengths mismatch:", i, len(ar[i]), len(ap[i]))
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
	d1 := fstU4(1, ar, srt) // median of four different sorts
	d2 := fstU4(2, ar, srt)
	d3 := fstU4(3, ar, srt)
	d1 = medur(fstU4(4, ar, srt), d1, d2, d3)

	compareU4(ar, ap)
	return printSec(tn, d1)
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
	d1 := fstF4(5, ar, srt) // median of four different sorts
	d2 := fstF4(6, ar, srt)
	d3 := fstF4(7, ar, srt)
	d1 = medur(fstF4(8, ar, srt), d1, d2, d3)

	compareU4(F4toU4(&ar), F4toU4(&ap))
	return printSec(tn, d1)
}

// median fst & compare for string
func mfcS(tn string, srt func([]string), ar, ap []uint32) float64 {
	d1 := fstS(9, ar, srt) // median of four different sorts
	d2 := fstS(10, ar, srt)
	d3 := fstS(11, ar, srt)
	d1 = medur(fstS(12, ar, srt), d1, d2, d3)

	if len(ap) > 0 {
		as, ar := implantS(ar, false)
		aq, ap := implantS(ap, false)
		compareS(as, aq)
		compareU4(ar, ap)
	}
	return printSec(tn, d1)
}

// median fst & compare for []byte
func mfcB(tn string, srt func([][]byte), ar, ap []uint32) float64 {
	d1 := fstB(13, ar, srt) // median of four different sorts
	d2 := fstB(14, ar, srt)
	d3 := fstB(15, ar, srt)
	d1 = medur(fstB(16, ar, srt), d1, d2, d3)

	if len(ap) > 0 {
		as, ar := implantB(ar, false)
		aq, ap := implantB(ap, false)
		compareB(as, aq)
		compareU4(ar, ap)
	}
	return printSec(tn, d1)
}

// median fst & compare for string (SortLen)
func mfcLenS(tn string, srt func([]string), ar, ap []uint32) float64 {
	d1 := fstLenS(17, ar, srt) // median of four different sorts
	d2 := fstLenS(18, ar, srt)
	d3 := fstLenS(19, ar, srt)
	d1 = medur(fstLenS(20, ar, srt), d1, d2, d3)

	if len(ap) > 0 {
		as := implantLenS(0, ar, false)
		aq := implantLenS(0, ap, false)
		compareLenS(as, aq)
	}
	return printSec(tn, d1)
}

// median fst & compare for []byte (SortLen)
func mfcLenB(tn string, srt func([][]byte), ar, ap []uint32) float64 {
	d1 := fstLenB(21, ar, srt) // median of four different sorts
	d2 := fstLenB(22, ar, srt)
	d3 := fstLenB(23, ar, srt)
	d1 = medur(fstLenB(24, ar, srt), d1, d2, d3)

	if len(ap) > 0 {
		as := implantLenB(0, ar, false)
		aq := implantLenB(0, ap, false)
		compareLenB(as, aq)
	}
	return printSec(tn, d1)
}

// return sum of sortU4() times for 1..4 goroutines
// compare with ap and among themselves
func sumtU4(ar, ap []uint32) float64 {
	s := .0
	for MaxGor = 1; MaxGor < 5; MaxGor++ {
		s += mfcU4(fmt.Sprintf("sorty-%d", MaxGor), sortU4, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of sortF4() times for 1..4 goroutines
// compare with ap and among themselves
func sumtF4(ar, ap []float32) float64 {
	s := .0
	for MaxGor = 1; MaxGor < 5; MaxGor++ {
		s += mfcF4(fmt.Sprintf("sorty-%d", MaxGor), sortF4, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of sortS() times for 1..4 goroutines
// compare with ap and among themselves
func sumtS(ar, ap []uint32) float64 {
	s := .0
	for MaxGor = 1; MaxGor < 5; MaxGor++ {
		s += mfcS(fmt.Sprintf("sorty-%d", MaxGor), sortS, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of sortB() times for 1..4 goroutines
// compare with ap and among themselves
func sumtB(ar, ap []uint32) float64 {
	s := .0
	for MaxGor = 1; MaxGor < 5; MaxGor++ {
		s += mfcB(fmt.Sprintf("sorty-%d", MaxGor), sortB, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of SortLen([]string) times for 1..4 goroutines
// compare with ap and among themselves
func sumtLenS(ar, ap []uint32) float64 {
	s := .0
	for MaxGor = 1; MaxGor < 5; MaxGor++ {
		s += mfcLenS(fmt.Sprintf("sorty-%d", MaxGor), func(al []string) { SortLen(al) }, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of SortLen([][]byte) times for 1..4 goroutines
// compare with ap and among themselves
func sumtLenB(ar, ap []uint32) float64 {
	s := .0
	for MaxGor = 1; MaxGor < 5; MaxGor++ {
		s += mfcLenB(fmt.Sprintf("sorty-%d", MaxGor), func(al [][]byte) { SortLen(al) }, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// sort uint32 slice with Sort()
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

// return sum of sort3i() times for 1..4 goroutines
// compare with ap and among themselves
func sumtLswU4(ar, ap []uint32) float64 {
	s := .0
	for MaxGor = 1; MaxGor < 5; MaxGor++ {
		s += mfcU4(fmt.Sprintf("sortyLsw-%d", MaxGor), sort3i, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// sort float32 slice with Sort()
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
	for MaxGor = 1; MaxGor < 5; MaxGor++ {
		s += mfcF4(fmt.Sprintf("sortyLsw-%d", MaxGor), sort3f, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// sort string slice with Sort()
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
	for MaxGor = 1; MaxGor < 5; MaxGor++ {
		s += mfcS(fmt.Sprintf("sortyLsw-%d", MaxGor), sort3s, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}
