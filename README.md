## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![go.dev ref](/.github/godev.svg)](https://pkg.go.dev/github.com/jfcg/sorty)
Type-specific, fast, concurrent/parallel sorting library.

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation (with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine) and does not require extra memory. Call corresponding Sort\*() to concurrently sort your slice (in ascending order) or collection. For example:
```
sorty.SortS(string_slice) // native slice
sorty.Sort(n, lesswap)    // lesswap() function based
```
If you have a pair of `Less()` and `Swap()`, then you can trivially write your [lesswap()](https://pkg.go.dev/github.com/jfcg/sorty#Sort) and sort your collection concurrently.
Also, `lesswap()` operates faster than `sort.Interface` on generic collections.

For each Sort\*() call, sorty uses up to `Mxg` (3 by default, including caller) concurrent goroutines and one channel. They are created **only when necessary**. Its `Mli/Hmli/Mlr` parameters are tuned to get the best performance, see below.
Also, sorty uses [semantic](https://semver.org) versioning.

### Benchmarks
Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go version `1.15.5` on:

Machine|CPU|OS|Kernel
:---|:---|:---|:---
R |Ryzen 1600   |Manjaro     |5.4.74
X |Xeon Haswell |Ubuntu 20.04|5.4.0-1029-gcp
i5|Core i5 4210M|Manjaro     |5.4.74

Sorting random uint32 array (in seconds):

Library(-Mxg)|R|X|i5
:---|:---:|:---:|:---:
sort.Slice|16.03|17.94|16.10
sortutil  | 2.33| 2.93| 3.85
zermelo   | 2.01| 1.76| 1.09
sorty-2   | 3.32| 3.76| 3.21
sorty-3   | 2.30| 2.60| 2.66
sorty-4   | 1.77| 1.95| 2.30
sortyLsw-2| 7.23| 8.21| 7.66
sortyLsw-3| 4.92| 5.20| 6.02
sortyLsw-4| 4.14| 4.48| 5.46

Sorting random float32 array (in seconds):

Library(-Mxg)|R|X|i5
:---|:---:|:---:|:---:
sort.Slice|17.41|18.95|17.20
sortutil  | 2.90| 3.15| 4.47
zermelo   | 4.63| 3.98| 3.17
sorty-2   | 4.00| 4.19| 3.64
sorty-3   | 2.79| 2.72| 2.85
sorty-4   | 2.33| 2.36| 2.54
sortyLsw-2| 7.95| 8.59| 7.97
sortyLsw-3| 5.39| 5.84| 6.27
sortyLsw-4| 4.47| 4.48| 5.62

Sorting random string array (in seconds):

Library(-Mxg)|R|X|i5
:---|:---:|:---:|:---:
sort.Slice| 8.34| 9.48| 8.44
sortutil  | 1.65| 2.03| 2.27
radix     | 4.29| 4.40| 3.50
sorty-2   | 2.90| 3.55| 3.39
sorty-3   | 1.99| 2.48| 3.07
sorty-4   | 1.76| 1.96| 2.83
sortyLsw-2| 3.76| 4.55| 4.36
sortyLsw-3| 2.57| 2.94| 3.99
sortyLsw-4| 2.24| 2.50| 3.75

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
