/*	Copyright (c) 2021, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "reflect"

// IsSortedLen returns 0 if ar is sorted 'by length' in ascending order, otherwise
// it returns i > 0 with len(ar[i]) < len(ar[i-1]). ar's (underlying) type can be
//
//	[]string, [][]T // for any type T
//
// otherwise it panics.
//
//go:nosplit
func IsSortedLen(ar any) int {
	ptr, len, kind := extractSK(ar)
	switch {
	case kind == reflect.String:
		return isSortedHL(toSlc[string](ptr, len))
	case kind >= sliceBias:
		return isSortedHL(toSlc[[]struct{}](ptr, len))
	}
	panic("sorty: IsSortedLen: invalid input type")
}

// SortLen concurrently sorts ar 'by length' in ascending order. ar's (underlying)
// type can be
//
//	[]string, [][]T // for any type T
//
// otherwise it panics.
//
//go:nosplit
func SortLen(ar any) {
	ptr, len, kind := extractSK(ar)
	switch {
	case kind == reflect.String:
		sortHL(toSlc[string](ptr, len))
	case kind >= sliceBias:
		sortHL(toSlc[[]struct{}](ptr, len))
	default:
		panic("sorty: SortLen: invalid input type")
	}
}
