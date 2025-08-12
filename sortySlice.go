/*	Copyright (c) 2021, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"reflect"

	sb "github.com/jfcg/sixb/v2"
)

// IsSortedSlice returns 0 if ar is sorted in ascending order, otherwise
// it returns i > 0 with ar[i] < ar[i-1]. ar's (underlying) type can be
//
//	[]int, []int32, []int64, []uint, []uint32, []uint64,
//	[]uintptr, []float32, []float64, []string, [][]byte,
//	[]unsafe.Pointer, []*T // for any type T
//
// otherwise it panics.
//
//go:nosplit
func IsSortedSlice(ar any) int {
	slc, kind := extractSK(ar)
	switch kind {
	case reflect.Int32:
		return isSortedO(sb.Cast[int32](slc))
	case reflect.Int64:
		return isSortedO(sb.Cast[int64](slc))
	case reflect.Uint32:
		return isSortedO(sb.Cast[uint32](slc))
	case reflect.Uint64:
		return isSortedO(sb.Cast[uint64](slc))
	case reflect.Float32:
		return isSortedF(sb.Cast[float32](slc))
	case reflect.Float64:
		return isSortedF(sb.Cast[float64](slc))
	case sliceBias + reflect.Uint8: // [][]byte
		return isSortedB(sb.Cast[[]byte](slc))
	case reflect.String:
		return isSortedO(sb.Cast[string](slc))
	}
	panic("sorty: IsSortedSlice: invalid input type")
}

// SortSlice concurrently sorts ar in ascending order. ar's (underlying) type can be
//
//	[]int, []int32, []int64, []uint, []uint32, []uint64,
//	[]uintptr, []float32, []float64, []string, [][]byte,
//	[]unsafe.Pointer, []*T // for any type T
//
// otherwise it panics.
func SortSlice(ar any) {
	slc, kind := extractSK(ar)
	switch kind {
	case reflect.Int32:
		sortI(sb.Cast[int32](slc))
	case reflect.Int64:
		sortI(sb.Cast[int64](slc))
	case reflect.Uint32:
		sortI(sb.Cast[uint32](slc))
	case reflect.Uint64:
		sortI(sb.Cast[uint64](slc))
	case reflect.Float32:
		sortF(sb.Cast[float32](slc))
	case reflect.Float64:
		sortF(sb.Cast[float64](slc))
	case sliceBias + reflect.Uint8: // [][]byte
		sortB(sb.Cast[[]byte](slc))
	case reflect.String:
		sortS(sb.Cast[string](slc))
	default:
		panic("sorty: SortSlice: invalid input type")
	}
}
