/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// Package sorty is a type-specific, fast, efficient, concurrent / parallel sorting
// library. It is an innovative [QuickSort] implementation, hence in-place and does not
// require extra memory. You can call:
//
//	import "github.com/jfcg/sorty/v2"
//
//	sorty.SortSlice(native_slice) // []int, []float64, []string, []*T etc. in ascending order
//	sorty.SortLen(len_slice)      // []string or [][]T 'by length' in ascending order
//	sorty.Sort(n, lesswap)        // lesswap() based
//
// [QuickSort]: https://en.wikipedia.org/wiki/Quicksort
package sorty

import (
	"reflect"
	"unsafe"

	"github.com/jfcg/sixb"
)

// MaxGor is the maximum number of goroutines (including caller) that can be
// concurrently used for sorting per Sort*() call. MaxGor can be changed live, even
// during ongoing Sort*() calls. MaxGor ≤ 1 (or a short input) yields single-goroutine
// sorting: sorty will not create any goroutines or channel.
var MaxGor uint32 = 3

func init() {
	if !(4097 > MaxGor && MaxGor > 0 && MaxLenRec > 2*MaxLenIns &&
		MaxLenIns > MaxLenInsFC && MaxLenInsFC > 2*nsShort) {
		panic("sorty: check your MaxGor/MaxLen* values")
	}
}

// Search returns lowest integer k in [0,n) where fn(k) is true, assuming:
//
//	fn(k) implies fn(k+1)
//
// If there is no such k, it returns n. It can be used to locate an element
// in a sorted collection.
func Search(n int, fn func(int) bool) int {
	l, h := 0, n

	for l < h {
		m := sixb.MeanI(l, h)

		if fn(m) {
			h = m
		} else {
			l = m + 1
		}
	}
	return l
}

// synchronization variables for [g]long*()
type syncVar struct {
	ngr  uint32   // number of sorting goroutines
	done chan int // end signal
}

// gorFull returns true if goroutine quota is full
//
//go:norace
func gorFull(sv *syncVar) bool {
	mg := MaxGor
	return sv.ngr >= mg
}

const (
	// #samples in pivot selection for
	nsShort = 4 // short range
	nsLong  = 6 // long range
	nsConc  = 8 // dual range
)

// Given n ≥ 2 and slice length ≥ 2n, select n equidistant samples
// from slice that minimizes max distance to non-selected members.
func minMaxSample(slen, n uint) (first, step, last uint) {
	step = slen / n // ≥ 2
	n--
	span := n * step
	tail := slen - span // 1 + #members in both tails
	if tail > n && tail>>1 > (step+1)>>1 {
		step++
		span += n
		tail -= n
	}
	first = tail >> 1 // larger tail
	last = first + span
	return
}

func insertionI(slc []int) {
	if unsafe.Sizeof(int(0)) == 8 {
		insertionI8(*(*[]int64)(unsafe.Pointer(&slc)))
	} else {
		insertionI4(*(*[]int32)(unsafe.Pointer(&slc)))
	}
}

const sliceBias reflect.Kind = 100

// extracts slice and element kind from ar
func extractSK(ar interface{}) (slc sixb.Slice, kind reflect.Kind) {
	tipe := reflect.TypeOf(ar)
	if tipe.Kind() != reflect.Slice {
		return
	}
	tipe = tipe.Elem()
	kind = tipe.Kind()

	switch kind {
	// map int/uint/pointer types to hardware type
	case reflect.Uintptr, reflect.Pointer, reflect.UnsafePointer:
		kind = reflect.Uint32 + reflect.Kind(unsafe.Sizeof(uintptr(0))>>3)
	case reflect.Uint:
		kind = reflect.Uint32 + reflect.Kind(unsafe.Sizeof(uint(0))>>3)
	case reflect.Int:
		kind = reflect.Int32 + reflect.Kind(unsafe.Sizeof(int(0))>>3)
	// map []T to sliceBias + Kind(T)
	case reflect.Slice:
		kind = sliceBias + tipe.Elem().Kind()
	// other recognized types
	case reflect.Int32, reflect.Int64, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
	default:
		kind = reflect.Invalid
		return
	}

	v := reflect.ValueOf(ar)
	p, l := v.Pointer(), v.Len()
	slc = sixb.Slice{Data: unsafe.Pointer(p), Len: l, Cap: l}
	return
}
