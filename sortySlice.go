/*	Copyright (c) 2021, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "reflect"

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
	ptr, len, kind := extractSK(ar)
	switch kind {
	case reflect.Int32:
		return isSortedO(toSlc[int32](ptr, len))
	case reflect.Int64:
		return isSortedO(toSlc[int64](ptr, len))
	case reflect.Uint32:
		return isSortedO(toSlc[uint32](ptr, len))
	case reflect.Uint64:
		return isSortedO(toSlc[uint64](ptr, len))
	case reflect.Float32:
		return isSortedF(toSlc[float32](ptr, len))
	case reflect.Float64:
		return isSortedF(toSlc[float64](ptr, len))
	case sliceBias + reflect.Uint8: // [][]byte
		return isSortedB(toSlc[[]byte](ptr, len))
	case reflect.String:
		return isSortedO(toSlc[string](ptr, len))
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
	ptr, len, kind := extractSK(ar)
	switch kind {
	case reflect.Int32:
		sortI(toSlc[int32](ptr, len))
	case reflect.Int64:
		sortI(toSlc[int64](ptr, len))
	case reflect.Uint32:
		sortI(toSlc[uint32](ptr, len))
	case reflect.Uint64:
		sortI(toSlc[uint64](ptr, len))
	case reflect.Float32:
		sortF(toSlc[float32](ptr, len))
	case reflect.Float64:
		sortF(toSlc[float64](ptr, len))
	case sliceBias + reflect.Uint8: // [][]byte
		sortB(toSlc[[]byte](ptr, len))
	case reflect.String:
		sortS(toSlc[string](ptr, len))
	default:
		panic("sorty: SortSlice: invalid input type")
	}
}
