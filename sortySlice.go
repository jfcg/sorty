/*	Copyright (c) 2021, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"reflect"
	"unsafe"

	"github.com/jfcg/sixb"
)

// extracts slice and element kind from ar
func extractSlc(ar interface{}) (slc sixb.Slice, k reflect.Kind) {
	t := reflect.TypeOf(ar)
	if t.Kind() != reflect.Slice {
		return
	}
	t = t.Elem()
	k = t.Kind()

	// map int,uint,uintptr to hw type
	if k == reflect.Int {
		k = reflect.Int32 + reflect.Kind(unsafe.Sizeof(int(0))>>3)
	} else if k == reflect.Uint || k == reflect.Uintptr {
		k = reflect.Uint32 + reflect.Kind(unsafe.Sizeof(uint(0))>>3)
	}

	if k == reflect.Int32 || k == reflect.Int64 || k == reflect.String ||
		reflect.Uint32 <= k && k <= reflect.Float64 || // uint32/64, float32/64
		k == reflect.Slice && t.Elem().Kind() == reflect.Uint8 { // [][]byte

		v := reflect.ValueOf(ar)
		p, l := v.Pointer(), v.Len()
		slc = sixb.Slice{Data: unsafe.Pointer(p), Len: l, Cap: l}
		return
	}
	k = reflect.Invalid
	return
}

// IsSortedSlice returns 0 if ar is sorted in ascending order, otherwise it returns
// i > 0 with ar[i] < ar[i-1]. ar's (underlying) type can be []int, []int32, []int64,
// []uint, []uint32, []uint64, []uintptr, []float32, []float64, []string, [][]byte,
// otherwise it panics.
func IsSortedSlice(ar interface{}) int {
	slc, k := extractSlc(ar)
	switch k {
	case reflect.Int32:
		i := *(*[]int32)(unsafe.Pointer(&slc))
		return isSortedI4(i)
	case reflect.Int64:
		i := *(*[]int64)(unsafe.Pointer(&slc))
		return isSortedI8(i)
	case reflect.Uint32:
		u := *(*[]uint32)(unsafe.Pointer(&slc))
		return isSortedU4(u)
	case reflect.Uint64:
		u := *(*[]uint64)(unsafe.Pointer(&slc))
		return isSortedU8(u)
	case reflect.Float32:
		f := *(*[]float32)(unsafe.Pointer(&slc))
		return isSortedF4(f)
	case reflect.Float64:
		f := *(*[]float64)(unsafe.Pointer(&slc))
		return isSortedF8(f)
	case reflect.Slice:
		b := *(*[][]byte)(unsafe.Pointer(&slc))
		return isSortedB(b)
	case reflect.String:
		s := *(*[]string)(unsafe.Pointer(&slc))
		return isSortedS(s)
	}
	panic("sorty: IsSortedSlice: invalid input type")
}

// SortSlice concurrently sorts ar in ascending order. ar's (underlying) type can be
// []int, []int32, []int64, []uint, []uint32, []uint64, []uintptr, []float32, []float64,
// []string, [][]byte, otherwise it panics.
func SortSlice(ar interface{}) {
	slc, k := extractSlc(ar)
	switch k {
	case reflect.Int32:
		i := *(*[]int32)(unsafe.Pointer(&slc))
		sortI4(i)
	case reflect.Int64:
		i := *(*[]int64)(unsafe.Pointer(&slc))
		sortI8(i)
	case reflect.Uint32:
		u := *(*[]uint32)(unsafe.Pointer(&slc))
		sortU4(u)
	case reflect.Uint64:
		u := *(*[]uint64)(unsafe.Pointer(&slc))
		sortU8(u)
	case reflect.Float32:
		f := *(*[]float32)(unsafe.Pointer(&slc))
		sortF4(f)
	case reflect.Float64:
		f := *(*[]float64)(unsafe.Pointer(&slc))
		sortF8(f)
	case reflect.Slice:
		b := *(*[][]byte)(unsafe.Pointer(&slc))
		sortB(b)
	case reflect.String:
		s := *(*[]string)(unsafe.Pointer(&slc))
		sortS(s)
	default:
		panic("sorty: SortSlice: invalid input type")
	}
}
