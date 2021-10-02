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

func extract(ar interface{}) (slc sixb.Slice, r int) {
	t := reflect.TypeOf(ar)
	if t.Kind() != reflect.Slice {
		return
	}
	k := t.Elem().Kind()
	if k == reflect.String {
		r++
	} else if k == reflect.Slice {
		r--
	} else {
		return
	}

	v := reflect.ValueOf(ar)
	p, l := v.Pointer(), v.Len()
	slc = sixb.Slice{Data: unsafe.Pointer(p), Len: l, Cap: l}
	return
}

// IsSortedLen returns 0 if ar is sorted 'by length' in ascending order, otherwise
// it returns i > 0 with len(ar[i]) < len(ar[i-1]). ar's (underlying) type can be
// []string or [][]T (for any type T), otherwise it panics.
func IsSortedLen(ar interface{}) int {
	slc, r := extract(ar)
	if r == 0 {
		panic("sorty: IsSortedLen: invalid input type")
	}
	if r > 0 {
		s := *(*[]string)(unsafe.Pointer(&slc))
		for i := len(s) - 1; i > 0; i-- {
			if len(s[i]) < len(s[i-1]) {
				return i
			}
		}
		return 0
	}
	b := *(*[][]byte)(unsafe.Pointer(&slc))
	for i := len(b) - 1; i > 0; i-- {
		if len(b[i]) < len(b[i-1]) {
			return i
		}
	}
	return 0
}

// SortLen concurrently sorts ar 'by length' in ascending order. ar's (underlying)
// type can be []string or [][]T (for any type T), otherwise it panics.
func SortLen(ar interface{}) {
	slc, r := extract(ar)
	if r == 0 {
		panic("sorty: SortLen: invalid input type")
	}
	if r > 0 {
		s := *(*[]string)(unsafe.Pointer(&slc))
		sortLenS(s)
		return
	}
	b := *(*[][]byte)(unsafe.Pointer(&slc))
	sortLenB(b)
}
