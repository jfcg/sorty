/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// Package sorty provides type-specific concurrent / parallel sorting functionality.
//
// sorty is an in-place QuickSort implementation and does not require extra memory.
// Call corresponding Sort*() to concurrently sort your slice (in ascending order)
// or collection. For example:
//  sorty.SortS(string_slice) // native slice
//  sorty.Sort(n, lesswap)    // lesswap() function based
package sorty

var (
	// Mxg is the maximum number of goroutines used for sorting per Sort*() call.
	Mxg uint32 = 3

	// Mli is the maximum array length for insertion sort in
	// Sort*() except SortS() and Sort().
	Mli = 108
	// Hmli is the maximum array length for insertion sort in SortS() and Sort().
	Hmli = (Mli + 1) / 3

	// Mlr is the maximum array length for recursion when there is available goroutines.
	// So Mlr+1 is the minimum array length for new sorting goroutines.
	Mlr = 496
)

func init() {
	if !(4097 > Mxg && Mxg > 0 && Mlr > 2*Mli && Mli > Hmli && Hmli > 15) {
		panic("sorty: check your Mxg/Mli/Hmli/Mlr values")
	}
}

// mid-point
func mid(l, h int) int {
	return int(uint(l+h) >> 1)
}

// Search returns lowest integer k in [0,n) where fn(k) is true, assuming:
//  fn(k) => fn(k+1)
// If there is no such k, it returns n. It can be used to locate an element
// in a sorted array or collection.
func Search(n int, fn func(int) bool) int {
	l, h := 0, n

	for l < h {
		m := mid(l, h)

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
