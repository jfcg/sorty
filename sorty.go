// Type-specific concurrent sorting library
//
// Assign your slice to Ar* and call Sort*() to sort ascending. For example:
//  sorty.ArS = your_string_slice
//  sorty.SortS()
// There is no limit on number of goroutines to create, though sorty does it sparingly.
package sorty

// Minimum array size for Quick Sort*()
const S = 25
