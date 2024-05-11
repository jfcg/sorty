/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"reflect"
	"sort"
	"testing"
	"time"
	"unsafe"

	"github.com/jfcg/rng"
	"github.com/jfcg/sixb"
)

const (
	// will have many equal & distinct elements in buffers
	bufFull = 10_000_000
	bufHalf = bufFull / 2

	maxMaxGor = 3
)

var (
	// A & B buffers will hold all slices to sort
	aaBuf = make([]uint32, bufFull)
	bbBuf = make([]uint32, bufFull)

	// source buffer
	srcBuf = make([]uint32, bufFull)

	tsPtr *testing.T
)

// fill source buffer with random bytes
func fillSrc() {
	rng.Fill(sixb.U4toB(srcBuf))
}

// copy, prepare, sort, test
func copyPrepSortTest(buf []uint32, prepare func([]uint32) any,
	srf func(any)) (time.Duration, any) {

	// copy from srcBuf into buf
	if p := &buf[0]; p == &aaBuf[0] || p == &bbBuf[0] {
		copy(buf, srcBuf)
	} else {
		copy(buf, srcBuf[bufHalf:]) // initialize from second half
	}

	// prepare input if given
	// its output type determines the sort type
	var ar any
	if prepare != nil {
		ar = prepare(buf)
	} else {
		ar = buf
	}

	// measure duration of sorting
	now := time.Now()
	srf(ar)
	dur := time.Since(now)

	isSorted := IsSortedSlice // by value
	if !isValueSort(srf) {
		isSorted = IsSortedLen // by length
	}
	// check if result is sorted
	if isSorted(ar) != 0 {
		_, kind := extractSK(ar)
		tsPtr.Fatal("not sorted, kind:", kind)
	}
	return dur, ar
}

func isValueSort(srf func(any)) bool {
	sPtr := reflect.ValueOf(srf).Pointer()
	return sPtr == stdSortPtr || sPtr == sortSlcPtr || sPtr == sortLswPtr
}

var (
	stdSortPtr = reflect.ValueOf(stdSort).Pointer()   // standard sort.Slice
	sortSlcPtr = reflect.ValueOf(SortSlice).Pointer() // sorty
	sortLswPtr = reflect.ValueOf(sortLsw).Pointer()   // sorty
)

func stdSort(ar any) {
	slc, kind := extractSK(ar)

	switch kind {
	case reflect.Float32:
		buf := *(*[]float32)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool {
			x, y := buf[i], buf[k]
			return x < y || NaNoption == NaNlarge && x == x && y != y ||
				NaNoption == NaNsmall && x != x && y == y
		})
	case reflect.Float64:
		buf := *(*[]float64)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool {
			x, y := buf[i], buf[k]
			return x < y || NaNoption == NaNlarge && x == x && y != y ||
				NaNoption == NaNsmall && x != x && y == y
		})
	case reflect.Int32:
		buf := *(*[]int32)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case reflect.Int64:
		buf := *(*[]int64)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case reflect.Uint32:
		buf := *(*[]uint32)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case reflect.Uint64:
		buf := *(*[]uint64)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case reflect.String:
		buf := *(*[]string)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case sliceBias + reflect.Uint8:
		buf := *(*[][]byte)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool { return sixb.BtoS(buf[i]) < sixb.BtoS(buf[k]) })
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}
}

//go:nosplit
func stdSortLen(ar any) {
	slc, kind := extractSK(ar)

	switch {
	case kind == reflect.String:
		buf := *(*[]string)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool { return len(buf[i]) < len(buf[k]) })
	case kind >= sliceBias:
		buf := *(*[][]byte)(unsafe.Pointer(&slc))
		sort.Slice(buf, func(i, k int) bool { return len(buf[i]) < len(buf[k]) })
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}
}

func basicCheck(ar, ap any) (slc1, slc2 sixb.Slice, kind reflect.Kind) {
	slc1, kind = extractSK(ar)
	slc2, kind2 := extractSK(ap)
	if kind != kind2 {
		tsPtr.Fatal("different kinds:", kind, kind2)
	}
	if slc1.Len != slc2.Len {
		tsPtr.Fatal("length mismatch:", kind, slc1.Len, slc2.Len)
	}
	if slc1.Data == slc2.Data {
		tsPtr.Fatal("same slice data:", kind, slc1.Data)
	}
	return
}

func compare(ar, ap any) { // by value
	slc1, slc2, kind := basicCheck(ar, ap)
	var buf1, buf2 []uint64

	switch kind {
	case reflect.String:
		buf1 := *(*[]string)(unsafe.Pointer(&slc1))
		buf2 := *(*[]string)(unsafe.Pointer(&slc2))
		for i := len(buf1) - 1; i >= 0; i-- {
			if buf1[i] != buf2[i] {
				tsPtr.Fatal("values mismatch:", kind, i, buf1[i], buf2[i])
			}
		}
		return
	case sliceBias + reflect.Uint8:
		buf1 := *(*[][]byte)(unsafe.Pointer(&slc1))
		buf2 := *(*[][]byte)(unsafe.Pointer(&slc2))
		for i := len(buf1) - 1; i >= 0; i-- {
			if a, b := sixb.BtoS(buf1[i]), sixb.BtoS(buf2[i]); a != b {
				tsPtr.Fatal("values mismatch:", kind, i, a, b)
			}
		}
		return
	case reflect.Float32:
		buf1 := *(*[]float32)(unsafe.Pointer(&slc1))
		buf2 := *(*[]float32)(unsafe.Pointer(&slc2))
		for i := len(buf1) - 1; i >= 0; i-- {
			a, b := buf1[i], buf2[i]
			if a != b && (a == a || b == b) { // consider NaNs equal
				tsPtr.Fatal("values mismatch:", kind, i, a, b)
			}
		}
		return
	case reflect.Float64:
		buf1 := *(*[]float64)(unsafe.Pointer(&slc1))
		buf2 := *(*[]float64)(unsafe.Pointer(&slc2))
		for i := len(buf1) - 1; i >= 0; i-- {
			a, b := buf1[i], buf2[i]
			if a != b && (a == a || b == b) { // consider NaNs equal
				tsPtr.Fatal("values mismatch:", kind, i, a, b)
			}
		}
		return
	case reflect.Int32, reflect.Uint32:
		b1 := *(*[]uint32)(unsafe.Pointer(&slc1))
		b2 := *(*[]uint32)(unsafe.Pointer(&slc2))
		buf1 = sixb.U4toU8(b1)
		buf2 = sixb.U4toU8(b2)
	case reflect.Int64, reflect.Uint64:
		buf1 = *(*[]uint64)(unsafe.Pointer(&slc1))
		buf2 = *(*[]uint64)(unsafe.Pointer(&slc2))
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}

	for i := len(buf1) - 1; i >= 0; i-- {
		if buf1[i] != buf2[i] {
			tsPtr.Fatal("values mismatch:", kind, i, buf1[i], buf2[i])
		}
	}
}

func compareLen(ar, ap any) { // by length
	slc1, slc2, kind := basicCheck(ar, ap)

	switch {
	case kind == reflect.String:
		buf1 := *(*[]string)(unsafe.Pointer(&slc1))
		buf2 := *(*[]string)(unsafe.Pointer(&slc2))
		for i := len(buf1) - 1; i >= 0; i-- {
			if a, b := len(buf1[i]), len(buf2[i]); a != b {
				tsPtr.Fatal("len values mismatch:", kind, i, a, b)
			}
		}
	case kind >= sliceBias:
		buf1 := *(*[][]byte)(unsafe.Pointer(&slc1))
		buf2 := *(*[][]byte)(unsafe.Pointer(&slc2))
		for i := len(buf1) - 1; i >= 0; i-- {
			if a, b := len(buf1[i]), len(buf2[i]); a != b {
				tsPtr.Fatal("len values mismatch:", kind, i, a, b)
			}
		}
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
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
	return time.Duration(sixb.MeanI8(int64(b), int64(c)))
}

// Calculate median duration of four distinct calls with random inputs to srf().
// Optionally compare each result with standard sort.Slice.
// prepare()'s output type determines the sort type.
func medianCpstCompare(testName string, prepare func([]uint32) any,
	srf func(any), compStd bool) float64 {

	var std func(any)
	var cmp func(any, any)

	if compStd {
		if isValueSort(srf) {
			std = stdSort
			cmp = compare // by value
		} else {
			std = stdSortLen
			cmp = compareLen // by length
		}
	}

	dur := [4]time.Duration{}
	var ar any

	for i := 0; i < len(dur); i++ {
		fillSrc()
		dur[i], ar = copyPrepSortTest(aaBuf, prepare, srf)
		if compStd {
			_, ap := copyPrepSortTest(bbBuf, prepare, std)
			cmp(ar, ap)
		}
	}

	dur[0] = medur(dur[0], dur[1], dur[2], dur[3])

	return printSec(testName, dur[0])
}

var stNames = [4]string{"sorty-1", "sorty-2", "sorty-3", "sorty-4"}

// return sum of sortU4() durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurU4(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(stNames[MaxGor-1], nil, SortSlice, compStd)
	}
	return
}

func U4toF4(buf []uint32) any {
	return *(*[]float32)(unsafe.Pointer(&buf))
}

// return sum of sortF4() durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurF4(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(stNames[MaxGor-1], U4toF4, SortSlice, compStd)
	}
	return
}

// implant strings into buf
func implantS(buf []uint32) any {
	// string size is 4*t bytes
	t := sixb.StrSize >> 2

	// buf will hold n strings (headers followed by overlapping 12-bytes bodies)
	n := (len(buf) - 2) / (t + 1)

	t *= n // total string headers space
	ss := sixb.U4toStrs(buf[:t:t])

	for k := len(buf) - 2; n > 0; {
		n--
		k--
		ss[n].Data = unsafe.Pointer(&buf[k])
		ss[n].Len = 12
	}
	return *(*[]string)(unsafe.Pointer(&ss))
}

// return sum of sortS() durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurS(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(stNames[MaxGor-1], implantS, SortSlice, compStd)
	}
	return
}

// implant []byte's into buf
func implantB(buf []uint32) any {
	// []byte size is 4*t bytes
	t := sixb.SliceSize >> 2

	// buf will hold n []byte's (headers followed by overlapping 12-bytes bodies)
	n := (len(buf) - 2) / (t + 1)

	t *= n // total []byte headers space
	bs := sixb.U4toSlcs(buf[:t:t])

	for k := len(buf) - 2; n > 0; {
		n--
		k--
		bs[n].Data = unsafe.Pointer(&buf[k])
		bs[n].Len = 12
		bs[n].Cap = 12
	}
	return *(*[][]byte)(unsafe.Pointer(&bs))
}

// return sum of sortB() durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurB(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(stNames[MaxGor-1], implantB, SortSlice, compStd)
	}
	return
}

// implant strings into buf for SortLen
func implantLenS(buf []uint32) any {
	// string size is 4*t bytes
	t := sixb.StrSize >> 2

	// buf will hold n string headers
	n := len(buf) / t

	t *= n // total string headers space
	ss := sixb.U4toStrs(buf[:t:t])

	for L := 4*uint(len(buf)) + 1; n > 0; {
		n--
		// string bodies start at &buf[0] with random lengths up to 4*len(buf) bytes
		ss[n].Data = unsafe.Pointer(&buf[0])
		l := uint(ss[n].Len) // random number from srcBuf
		ss[n].Len = int(l % L)
	}
	return *(*[]string)(unsafe.Pointer(&ss))
}

// return sum of sortLenS() durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurLenS(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(stNames[MaxGor-1], implantLenS, SortLen, compStd)
	}
	return
}

// implant []byte's into buf
func implantLenB(buf []uint32) any {
	// []byte size is 4*t bytes
	t := sixb.SliceSize >> 2

	// buf will hold n []byte headers
	n := len(buf) / t

	t *= n // total []byte headers space
	bs := sixb.U4toSlcs(buf[:t:t])

	for L := 4*uint(len(buf)) + 1; n > 0; {
		n--
		// []byte bodies start at &buf[0] with random lengths up to 4*len(buf) bytes
		bs[n].Data = unsafe.Pointer(&buf[0])
		l := uint(bs[n].Len) % L // random number from srcBuf
		bs[n].Len = int(l)
		bs[n].Cap = int(l)
	}
	return *(*[][]byte)(unsafe.Pointer(&bs))
}

// return sum of sortLenB() durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurLenB(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(stNames[MaxGor-1], implantLenB, SortLen, compStd)
	}
	return
}

//go:nosplit
func sortLsw(ar any) {
	slc, kind := extractSK(ar)

	switch kind {
	case reflect.Uint32:
		buf := *(*[]uint32)(unsafe.Pointer(&slc))
		lsw := func(i, k, r, s int) bool {
			if buf[i] < buf[k] {
				if r != s {
					buf[r], buf[s] = buf[s], buf[r]
				}
				return true
			}
			return false
		}
		Sort(len(buf), lsw)
	case reflect.Float32:
		buf := *(*[]float32)(unsafe.Pointer(&slc))
		lsw := func(i, k, r, s int) bool {
			x, y := buf[i], buf[k]
			if x < y || NaNoption == NaNlarge && x == x && y != y ||
				NaNoption == NaNsmall && x != x && y == y {
				if r != s {
					buf[r], buf[s] = buf[s], buf[r]
				}
				return true
			}
			return false
		}
		Sort(len(buf), lsw)
	case reflect.String:
		buf := *(*[]string)(unsafe.Pointer(&slc))
		lsw := func(i, k, r, s int) bool {
			if buf[i] < buf[k] {
				if r != s {
					buf[r], buf[s] = buf[s], buf[r]
				}
				return true
			}
			return false
		}
		Sort(len(buf), lsw)
	case sliceBias + reflect.Uint8:
		buf := *(*[][]byte)(unsafe.Pointer(&slc))
		lsw := func(i, k, r, s int) bool {
			if sixb.BtoS(buf[i]) < sixb.BtoS(buf[k]) {
				if r != s {
					buf[r], buf[s] = buf[s], buf[r]
				}
				return true
			}
			return false
		}
		Sort(len(buf), lsw)
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}
}

var sLswNames = [4]string{"sortyLsw-1", "sortyLsw-2", "sortyLsw-3", "sortyLsw-4"}

// return sum of Sort([]uint32) durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurLswU4(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(sLswNames[MaxGor-1], nil, sortLsw, compStd)
	}
	return
}

// return sum of Sort([]float32) durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurLswF4(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(sLswNames[MaxGor-1], U4toF4, sortLsw, compStd)
	}
	return
}

// return sum of Sort([]string) durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurLswS(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(sLswNames[MaxGor-1], implantS, sortLsw, compStd)
	}
	return
}

// return sum of Sort([][]byte) durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurLswB(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(sLswNames[MaxGor-1], implantB, sortLsw, compStd)
	}
	return
}
