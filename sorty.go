// Package sorty provides type-specific concurrent sorting functionality
//
// sorty is an in-place QuickSort implementation and does not require extra memory. Call corresponding Sort*() to sort your slice in ascending order. For example:
//  sorty.SortS(string_slice, mx)
// mx is the maximum number of goroutines used for sorting simultaneously.
package sorty

// Mli is the maximum array length for insertion sort
var Mli = 96

// Mlr is the maximum array length for recursion when there is available goroutines.
// So Mlr+1 is the minimum array length for new sorting goroutines.
var Mlr = 385

func init() {
	if Mli < 8 || Mlr < 17 || Mlr <= 2*Mli {
		panic("sorty: check your Mli/Mlr values")
	}
}

// saturate to [2, 65535]
func sat(mx uint32) uint32 {
	if mx&^65535 != 0 {
		return 65535
	}
	if mx&^1 == 0 {
		return 2
	}
	return mx
}
