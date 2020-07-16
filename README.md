## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![go.dev ref](/.github/godev.svg)](https://pkg.go.dev/github.com/jfcg/sorty?tab=doc)
Type-specific, fast concurrent / parallel sorting library.

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation (with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine) and does not require extra memory. Call corresponding Sort\*() to concurrently sort your slice (in ascending order) or collection. For example:
```
sorty.SortS(string_slice) // native slice
sorty.Sort(n, lesswap)    // lesswap() function based
```
If you have a working `Less()` and `Swap()`, then you can trivially write your [lesswap()](https://pkg.go.dev/github.com/jfcg/sorty?tab=doc#Sort) and sort your collection concurrently.
`lesswap()` function is faster than `sort.Interface` as a way to access generic collections.

Mxg (3 by default) is the maximum number of goroutines used for sorting per Sort\*() call.
sorty uses [semantic](https://semver.org) versioning.

### Benchmarks
Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with:
```
Go version: 1.14.5
Machine R : Manjaro      on Ryzen 1600,   kernel 5.4.44
Machine X : Ubuntu 20.04 on Xeon Haswell, kernel 5.4.0-1019-gcp
```
Sorting uint32 array (in seconds):

Library(-Mxg)|R|X
:---|:---:|:---:
sort.Slice|16.21|20.01
sortutil  | 2.40| 3.04
zermelo   | 2.01| 1.79
sorty-2   | 3.19| 4.40
sorty-3   | 2.48| 3.16
sorty-4   | 2.08| 2.61
sortyLsw-2| 7.59| 9.23
sortyLsw-3| 4.95| 5.73
sortyLsw-4| 4.32| 5.00

Sorting float32 array (in seconds):

Library(-Mxg)|R|X
:---|:---:|:---:
sort.Slice|17.48|20.95
sortutil  | 2.95| 3.39
zermelo   | 4.53| 4.34
sorty-2   | 4.01| 4.76
sorty-3   | 3.00| 3.42
sorty-4   | 2.50| 2.78
sortyLsw-2| 8.18| 9.67
sortyLsw-3| 5.75| 6.33
sortyLsw-4| 4.29| 5.26

Sorting string array (in seconds):

Library(-Mxg)|R|X
:---|:---:|:---:
sort.Slice| 8.24|10.02
sortutil  | 1.65| 2.16
radix     | 4.25| 4.62
sorty-2   | 3.00| 3.96
sorty-3   | 2.33| 2.77
sorty-4   | 1.95| 2.27
sortyLsw-2| 3.89| 5.03
sortyLsw-3| 2.62| 3.33
sortyLsw-4| 2.34| 2.64

### Testing & Parameter Tuning
First, make sure everything is fine:
```
go test -short -timeout 1h
```
You can tune `Mli, Mlr` for your platform/cpu with (optimization flags):
```
go test -timeout 2h -gcflags '-B -wb=0 -smallframes' -ldflags '-s -w'
```
Now update `Mli, Mlr` in sorty.go and compare your tuned sorty with others:
```
go test -short -timeout 1h -gcflags '-B -wb=0 -smallframes' -ldflags '-s -w'
```
Remember to build sorty (and your functions like `SortObjAsc()`) with the same
optimization flags you used for tuning.

### Support
If you use sorty and like it, please support via ETH:`0x464B840ee70bBe7962b90bD727Aac172Fa8B9C15`
