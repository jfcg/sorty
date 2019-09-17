/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// Package sorty provides type-specific concurrent sorting functionality
//
// sorty is an in-place QuickSort implementation and does not require extra memory.
// Call corresponding Sort*() to concurrently sort your slice (in ascending order)
// or collection. For example:
//  sorty.SortS(string_slice)
//  sorty.Sort(col) // satisfies sort.Interface
package sorty

// Mxg is the maximum number of goroutines used for sorting per Sort*() call.
var Mxg uint32 = 3

// Mli is the maximum array length for insertion sort.
// Sort(Collection) uses 1/4 of this as its limit.
var Mli = 80

// Mlr is the maximum array length for recursion when there is available goroutines.
// So Mlr+1 is the minimum array length for new sorting goroutines.
var Mlr = 801

func init() {
	li2 := 2 * Mli
	if !(65536 > Mxg && Mxg > 1 && Mlr > li2 && li2 > 63) {
		panic("sorty: check your Mxg/Mli/Mlr values")
	}
}
