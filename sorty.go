/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// Package sorty provides type-specific concurrent / parallel sorting functionality
//
// sorty is an in-place QuickSort implementation and does not require extra memory.
// Call corresponding Sort*() to concurrently sort your slice (in ascending order)
// or collection. For example:
//  sorty.SortS(string_slice) // native slice
//  sorty.Sort(col)           // satisfies sort.Interface
//  sorty.Sort2(col2)         // satisfies sorty.Collection2
//  sorty.Sort3(n, lesswap)   // lesswap() function based
package sorty

// Mxg is the maximum number of goroutines used for sorting per Sort*() call.
var Mxg uint32 = 3

// Mli is the maximum array length for insertion sort.
// SortS() and Sort3() use 1/2 of this as their limits.
// Sort() and Sort2() use 1/4 of this as their limits.
var Mli = 96

// Mlr is the maximum array length for recursion when there is available goroutines.
// So Mlr+1 is the minimum array length for new sorting goroutines.
var Mlr = 401

func init() {
	li2 := 2 * Mli
	if !(65536 > Mxg && Mxg > 1 && Mlr > li2 && li2 > 63) {
		panic("sorty: check your Mxg/Mli/Mlr values")
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
