// Type-specific concurrent sorting library
//
// Call corresponding Sort\*() to sort your slice in ascending order. For example:
//  sorty.SortS(your_string_slice)
// A Sort\*() function should not be called by multiple goroutines at the same time. There is no limit on the number of goroutines to be created (could be many thousands depending on data), though sorty does it sparingly.
package sorty

// Minimum array size for Quick Sort*()
const S = 25
