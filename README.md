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
Go version: 1.15.4
Machine R : Manjaro      on Ryzen 1600,   kernel 5.4.74
Machine X : Ubuntu 20.04 on Xeon Haswell, kernel 5.4.0-1028-gcp
```
Sorting random uint32 array (in seconds):

Library(-Mxg)|R|X
:---|:---:|:---:
sort.Slice|16.10|17.23
sortutil  | 2.33| 2.84
zermelo   | 2.01| 1.70
sorty-2   | 3.33| 3.62
sorty-3   | 2.41| 2.52
sorty-4   | 2.00| 1.89
sortyLsw-2| 7.22| 7.82
sortyLsw-3| 4.84| 4.98
sortyLsw-4| 4.14| 4.28

Sorting random float32 array (in seconds):

Library(-Mxg)|R|X
:---|:---:|:---:
sort.Slice|17.41|18.13
sortutil  | 2.91| 3.11
zermelo   | 4.69| 3.81
sorty-2   | 3.98| 4.04
sorty-3   | 2.78| 2.62
sorty-4   | 2.30| 2.28
sortyLsw-2| 7.94| 8.22
sortyLsw-3| 5.56| 5.58
sortyLsw-4| 4.40| 4.32

Sorting random string array (in seconds):

Library(-Mxg)|R|X
:---|:---:|:---:
sort.Slice| 8.32| 8.94
sortutil  | 1.65| 1.93
radix     | 4.30| 4.19
sorty-2   | 2.90| 3.34
sorty-3   | 2.01| 2.34
sorty-4   | 1.76| 1.93
sortyLsw-2| 3.78| 4.32
sortyLsw-3| 2.58| 2.83
sortyLsw-4| 2.22| 2.39

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
