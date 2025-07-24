/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"bytes"
	"reflect"
	"slices"
	"sort"
	"testing"
	"time"
	"unsafe"

	"github.com/jfcg/rng"
	sb "github.com/jfcg/sixb/v2"
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
	rng.Fill(sb.Slice[byte](srcBuf))
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
	p := reflect.ValueOf(srf).Pointer()
	return p == stdSortPtr || p == stdSlicePtr || p == sortSlcPtr || p == sortLswPtr
}

var (
	stdSortPtr  = reflect.ValueOf(stdSort).Pointer()   // standard sort.Slice
	stdSlicePtr = reflect.ValueOf(stdSlice).Pointer()  // standard slices.Sort
	sortSlcPtr  = reflect.ValueOf(SortSlice).Pointer() // sorty
	sortLswPtr  = reflect.ValueOf(sortLsw).Pointer()   // sorty
)

func stdSort(ar any) {
	slc, kind := extractSK(ar)

	switch kind {
	case reflect.Float32:
		buf := sb.Cast[float32](slc)
		if NaNoption == NaNsmall {
			sort.Slice(buf, func(i, k int) bool {
				a, b := buf[i], buf[k]
				return a < b || a != a && b == b
			})
		} else {
			sort.Slice(buf, func(i, k int) bool {
				a, b := buf[i], buf[k]
				return a < b || a == a && b != b
			})
		}
	case reflect.Float64:
		buf := sb.Cast[float64](slc)
		if NaNoption == NaNsmall {
			sort.Slice(buf, func(i, k int) bool {
				a, b := buf[i], buf[k]
				return a < b || a != a && b == b
			})
		} else {
			sort.Slice(buf, func(i, k int) bool {
				a, b := buf[i], buf[k]
				return a < b || a == a && b != b
			})
		}
	case reflect.Int32:
		buf := sb.Cast[int32](slc)
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case reflect.Int64:
		buf := sb.Cast[int64](slc)
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case reflect.Uint32:
		buf := sb.Cast[uint32](slc)
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case reflect.Uint64:
		buf := sb.Cast[uint64](slc)
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case reflect.String:
		buf := sb.Cast[string](slc)
		sort.Slice(buf, func(i, k int) bool { return buf[i] < buf[k] })
	case sliceBias + reflect.Uint8:
		buf := sb.Cast[[]byte](slc)
		sort.Slice(buf, func(i, k int) bool {
			return sb.String(buf[i]) < sb.String(buf[k])
		})
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}
}

func stdSlice(ar any) {
	slc, kind := extractSK(ar)

	switch kind {
	case reflect.Float32:
		buf := sb.Cast[float32](slc)
		if NaNoption == NaNsmall {
			slices.Sort(buf)
		} else {
			slices.SortFunc(buf, func(a, b float32) int {
				if a < b || a == a && b != b {
					return -1
				}
				if a > b || a != a && b == b {
					return 1
				}
				return 0
			})
		}
	case reflect.Float64:
		buf := sb.Cast[float64](slc)
		if NaNoption == NaNsmall {
			slices.Sort(buf)
		} else {
			slices.SortFunc(buf, func(a, b float64) int {
				if a < b || a == a && b != b {
					return -1
				}
				if a > b || a != a && b == b {
					return 1
				}
				return 0
			})
		}
	case reflect.Int32:
		buf := sb.Cast[int32](slc)
		slices.Sort(buf)
	case reflect.Int64:
		buf := sb.Cast[int64](slc)
		slices.Sort(buf)
	case reflect.Uint32:
		buf := sb.Cast[uint32](slc)
		slices.Sort(buf)
	case reflect.Uint64:
		buf := sb.Cast[uint64](slc)
		slices.Sort(buf)
	case reflect.String:
		buf := sb.Cast[string](slc)
		slices.Sort(buf)
	case sliceBias + reflect.Uint8:
		buf := sb.Cast[[]byte](slc)
		slices.SortFunc(buf, bytes.Compare)
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}
}

//go:nosplit
func stdSortLen(ar any) {
	slc, kind := extractSK(ar)

	switch {
	case kind == reflect.String:
		buf := sb.Cast[string](slc)
		sort.Slice(buf, func(i, k int) bool { return len(buf[i]) < len(buf[k]) })
	case kind >= sliceBias:
		buf := sb.Cast[[]byte](slc)
		sort.Slice(buf, func(i, k int) bool { return len(buf[i]) < len(buf[k]) })
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}
}

//go:nosplit
func stdSliceLen(ar any) {
	slc, kind := extractSK(ar)

	switch {
	case kind == reflect.String:
		buf := sb.Cast[string](slc)
		slices.SortFunc(buf, func(a, b string) int { return len(a) - len(b) })
	case kind >= sliceBias:
		buf := sb.Cast[[]byte](slc)
		slices.SortFunc(buf, func(a, b []byte) int { return len(a) - len(b) })
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}
}

func basicCheck(ar, ap any) (slc1, slc2 sb.InSlice, kind reflect.Kind) {
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
		buf1 := sb.Cast[string](slc1)
		buf2 := sb.Cast[string](slc2)
		for i := len(buf1) - 1; i >= 0; i-- {
			if buf1[i] != buf2[i] {
				tsPtr.Fatal("values mismatch:", kind, i, buf1[i], buf2[i])
			}
		}
		return
	case sliceBias + reflect.Uint8:
		buf1 := sb.Cast[[]byte](slc1)
		buf2 := sb.Cast[[]byte](slc2)
		for i := len(buf1) - 1; i >= 0; i-- {
			if a, b := sb.String(buf1[i]), sb.String(buf2[i]); a != b {
				tsPtr.Fatal("values mismatch:", kind, i, a, b)
			}
		}
		return
	case reflect.Float32:
		buf1 := sb.Cast[float32](slc1)
		buf2 := sb.Cast[float32](slc2)
		for i := len(buf1) - 1; i >= 0; i-- {
			a, b := buf1[i], buf2[i]
			if a != b && (a == a || b == b) { // consider NaNs equal
				tsPtr.Fatal("values mismatch:", kind, i, a, b)
			}
		}
		return
	case reflect.Float64:
		buf1 := sb.Cast[float64](slc1)
		buf2 := sb.Cast[float64](slc2)
		for i := len(buf1) - 1; i >= 0; i-- {
			a, b := buf1[i], buf2[i]
			if a != b && (a == a || b == b) { // consider NaNs equal
				tsPtr.Fatal("values mismatch:", kind, i, a, b)
			}
		}
		return
	case reflect.Int32, reflect.Uint32:
		buf1 = sb.Slice[uint64](sb.Cast[uint32](slc1))
		buf2 = sb.Slice[uint64](sb.Cast[uint32](slc2))
	case reflect.Int64, reflect.Uint64:
		buf1 = sb.Cast[uint64](slc1)
		buf2 = sb.Cast[uint64](slc2)
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
		buf1 := sb.Cast[string](slc1)
		buf2 := sb.Cast[string](slc2)
		for i := len(buf1) - 1; i >= 0; i-- {
			if a, b := len(buf1[i]), len(buf2[i]); a != b {
				tsPtr.Fatal("len values mismatch:", kind, i, a, b)
			}
		}
	case kind >= sliceBias:
		buf1 := sb.Cast[[]byte](slc1)
		buf2 := sb.Cast[[]byte](slc2)
		for i := len(buf1) - 1; i >= 0; i-- {
			if a, b := len(buf1[i]), len(buf2[i]); a != b {
				tsPtr.Fatal("len values mismatch:", kind, i, a, b)
			}
		}
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}
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
			std = stdSlice
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

	dur[0] = sb.Median4(dur[0], dur[1], dur[2], dur[3])

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
	return sb.Slice[float32](buf)
}

// return sum of sortF4() durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurF4(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(stNames[MaxGor-1], U4toF4, SortSlice, compStd)
	}
	return
}

const (
	strSize = uint(unsafe.Sizeof(""))
	slcSize = uint(unsafe.Sizeof([]byte{}))
)

// implant strings into buf
func implantS(buf []uint32) any {
	// string size is 4*t bytes
	t := strSize >> 2

	// buf will hold n strings (headers followed by overlapping 12-bytes bodies)
	n := uint(len(buf)-2) / (t + 1)

	t *= n // total string headers space
	ss := sb.Slice[sb.InString](buf[:t:t])

	for k := len(buf) - 2; n > 0; {
		n--
		k--
		ss[n].Data = unsafe.Pointer(&buf[k])
		ss[n].Len = 12
	}
	return sb.Slice[string](ss)
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
	t := slcSize >> 2

	// buf will hold n []byte's (headers followed by overlapping 12-bytes bodies)
	n := uint(len(buf)-2) / (t + 1)

	t *= n // total []byte headers space
	bs := sb.Slice[sb.InSlice](buf[:t:t])

	for k := len(buf) - 2; n > 0; {
		n--
		k--
		bs[n].Data = unsafe.Pointer(&buf[k])
		bs[n].Len = 12
		bs[n].Cap = 12
	}
	return sb.Slice[[]byte](bs)
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
	t := strSize >> 2

	// buf will hold n string headers
	n := uint(len(buf)) / t

	t *= n // total string headers space
	ss := sb.Slice[sb.InString](buf[:t:t])

	for L := 4*uint(len(buf)) + 1; n > 0; {
		n--
		// string bodies start at &buf[0] with random lengths up to 4*len(buf) bytes
		ss[n].Data = unsafe.Pointer(&buf[0])
		// random number from srcBuf
		ss[n].Len %= L
	}
	return sb.Slice[string](ss)
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
	t := slcSize >> 2

	// buf will hold n []byte headers
	n := uint(len(buf)) / t

	t *= n // total []byte headers space
	bs := sb.Slice[sb.InSlice](buf[:t:t])

	for L := 4*uint(len(buf)) + 1; n > 0; {
		n--
		// []byte bodies start at &buf[0] with random lengths up to 4*len(buf) bytes
		bs[n].Data = unsafe.Pointer(&buf[0])
		l := bs[n].Len % L // random number from srcBuf
		bs[n].Len = l
		bs[n].Cap = l
	}
	return sb.Slice[[]byte](bs)
}

// return sum of sortLenB() durations for 1..maxMaxGor goroutines
// optionally compare with standard sort.Slice
func sumDurLenB(compStd bool) (sum float64) {
	for MaxGor = 1; MaxGor <= maxMaxGor; MaxGor++ {
		sum += medianCpstCompare(stNames[MaxGor-1], implantLenB, SortLen, compStd)
	}
	return
}

func sortLsw(ar any) {
	slc, kind := extractSK(ar)
	var lsw Lesswap

	switch kind {
	case reflect.Uint32:
		buf := sb.Cast[uint32](slc)
		lsw = func(i, k, r, s int) bool {
			if buf[i] < buf[k] {
				if r != s {
					buf[r], buf[s] = buf[s], buf[r]
				}
				return true
			}
			return false
		}
	case reflect.Float32:
		buf := sb.Cast[float32](slc)
		if NaNoption == NaNsmall {
			lsw = func(i, k, r, s int) bool {
				a, b := buf[i], buf[k]
				if a < b || a != a && b == b {
					if r != s {
						buf[r], buf[s] = buf[s], buf[r]
					}
					return true
				}
				return false
			}
		} else {
			lsw = func(i, k, r, s int) bool {
				a, b := buf[i], buf[k]
				if a < b || a == a && b != b {
					if r != s {
						buf[r], buf[s] = buf[s], buf[r]
					}
					return true
				}
				return false
			}
		}
	case reflect.String:
		buf := sb.Cast[string](slc)
		lsw = func(i, k, r, s int) bool {
			if buf[i] < buf[k] {
				if r != s {
					buf[r], buf[s] = buf[s], buf[r]
				}
				return true
			}
			return false
		}
	case sliceBias + reflect.Uint8:
		buf := sb.Cast[[]byte](slc)
		lsw = func(i, k, r, s int) bool {
			if sb.String(buf[i]) < sb.String(buf[k]) {
				if r != s {
					buf[r], buf[s] = buf[s], buf[r]
				}
				return true
			}
			return false
		}
	default:
		tsPtr.Fatal("unrecognized kind:", kind)
	}

	Sort(int(slc.Len), lsw)
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
