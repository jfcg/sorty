// Package sorty provides type-specific concurrent sorting functionality
//
// Call corresponding Sort*() to sort your slice in ascending order. For example:
//  sorty.SortS(your_string_slice, mx)
// A Sort*() function should not be called by multiple goroutines at the same time. mx is the maximum number of goroutines used for sorting simultaneously.
package sorty

// S is the minimum array size for Quick Sort*()
const S = 25
