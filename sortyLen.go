/*	Copyright (c) 2021, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"reflect"

	"github.com/jfcg/sixb/v2"
)

// IsSortedLen returns 0 if ar is sorted 'by length' in ascending order, otherwise
// it returns i > 0 with len(ar[i]) < len(ar[i-1]). ar's (underlying) type can be
//
//	[]string, [][]T // for any type T
//
// otherwise it panics.
//
//go:nosplit
func IsSortedLen(ar any) int {
	slc, kind := extractSK(ar)
	switch {
	case kind == reflect.String:
		return isSortedHL(sixb.Cast[string](slc))
	case kind >= sliceBias:
		return isSortedHL(sixb.Cast[[]struct{}](slc))
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
	slc, kind := extractSK(ar)
	switch {
	case kind == reflect.String:
		sortHL(sixb.Cast[string](slc))
	case kind >= sliceBias:
		sortHL(sixb.Cast[[]struct{}](slc))
	default:
		panic("sorty: SortLen: invalid input type")
	}
}
