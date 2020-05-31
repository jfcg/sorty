## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![go.dev ref](/.github/godev.svg)](https://pkg.go.dev/github.com/jfcg/sorty?tab=doc)
Type-specific, fast concurrent / parallel sorting library.

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation (with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine) and does not require extra memory. Call corresponding Sort\*() to concurrently sort your slice (in ascending order) or collection. For example:
```
sorty.SortS(string_slice) // native slice
sorty.Sort(n, lesswap)    // lesswap() function based
```
Mxg (3 by default) is the maximum number of goroutines used for sorting per Sort\*() call.
sorty uses [semantic](https://semver.org) versioning.

### Benchmarks
Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go 1.14.3.

Sorting uint32 array (in seconds):

Library(-Mxg)|Manjaro on Ryzen 1600|Ubuntu 20.04 on Xeon Haswell
:---|:---:|:---:
sort.Slice|16.03|20.24
sortutil  | 3.00| 3.34
zermelo   | 2.20| 1.75
sorty-2   | 3.29| 4.66
sorty-3   | 2.48| 3.66
sorty-4   | 2.05| 3.37
sortyLsw-2| 7.64| 9.21
sortyLsw-3| 5.00| 5.71
sortyLsw-4| 4.36| 4.98

Sorting float32 array (in seconds):

Library(-Mxg)|Manjaro on Ryzen 1600|Ubuntu 20.04 on Xeon Haswell
:---|:---:|:---:
sort.Slice|17.43|20.84
sortutil  | 3.01| 3.21
zermelo   | 4.69| 4.12
sorty-2   | 4.07| 4.76
sorty-3   | 3.03| 3.42
sorty-4   | 2.47| 2.77
sortyLsw-2| 8.31| 9.66
sortyLsw-3| 5.66| 6.33
sortyLsw-4| 4.61| 5.22

Sorting string array (in seconds):

Library(-Mxg)|Manjaro on Ryzen 1600|Ubuntu 20.04 on Xeon Haswell
:---|:---:|:---:
sort.Slice| 8.72| 9.52
sortutil  | 2.00| 2.01
radix     | 4.83| 4.41
sorty-2   | 3.24| 3.76
sorty-3   | 2.47| 2.62
sorty-4   | 2.07| 2.16
sortyLsw-2| 4.07| 4.83
sortyLsw-3| 2.88| 3.13
sortyLsw-4| 2.39| 2.54

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
