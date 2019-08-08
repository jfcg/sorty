// Package sorty provides type-specific concurrent sorting functionality
//
// sorty is an in-place QuickSort implementation and does not require extra memory. Call corresponding Sort*() to concurrently sort your slice in ascending order. For example:
//  sorty.SortS(string_slice)
package sorty

// Mxg is the maximum number of goroutines used for sorting per Sort*() call.
var Mxg uint32 = 3

// Mli is the maximum array length for insertion sort
var Mli = 64

// Mlr is the maximum array length for recursion when there is available goroutines.
// So Mlr+1 is the minimum array length for new sorting goroutines.
var Mlr = 321

func init() {
	li2 := 2 * Mli
	if !(65536 > Mxg && Mxg > 1 && Mlr > li2 && li2 > 15) {
		panic("sorty: check your Mxg/Mli/Mlr values")
	}
}
