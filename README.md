## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![go.dev ref](https://raw.githubusercontent.com/jfcg/.github/main/godev.svg)](https://pkg.go.dev/github.com/jfcg/sorty)
Type-specific, fast, efficient, concurrent/parallel sorting library.

sorty is a concurrent [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation (with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine). It is in-place and does not require extra memory. You can call corresponding `Sort*()` to rapidly sort your slice (in ascending order) or collection. For example:
```
sorty.SortS(string_slice) // native slice
sorty.Sort(n, lesswap)    // lesswap() function based
```
If you have a pair of `Less()` and `Swap()`, then you can trivially write your [`lesswap()`](https://pkg.go.dev/github.com/jfcg/sorty#Sort) and sort your generic collections using multiple cpu cores quickly.

sorty is stable, well tested and pretty careful with resources & performance:
- `lesswap()` operates faster than `sort.Interface` on generic collections.
- For each `Sort*()` call, sorty uses up to [`Mxg`](https://pkg.go.dev/github.com/jfcg/sorty#pkg-variables) (3 by default, including caller) concurrent goroutines and up to one channel.
- Goroutines and channel are created/used **only when necessary**.
- `Mxg=1` (or a short input) yields single-goroutine sorting, no goroutines or channel will be created.
- `Mxg` can be changed live, even during an ongoing `Sort*()` call.
- `Mli,Hmli,Mlr` parameters are tuned to get the best performance, see below.
- sorty API adheres to [semantic](https://semver.org) versioning.

### Benchmarks
Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go version `1.15.5` on:

Machine|CPU|OS|Kernel
:---:|:---|:---|:---
R |Ryzen 1600   |Manjaro     |5.4.80
X |Xeon Haswell |Ubuntu 20.04|5.4.0-1029-gcp
i5|Core i5 4210M|Manjaro     |5.4.80

Sorting random uint32 array (in seconds):

Library(-Mxg)|R|X|i5
:---|---:|---:|---:
sort.Slice|16.11|17.34|16.10
sortutil  | 2.35| 2.88| 3.87
zermelo   | 2.00| 1.74| 1.09
sorty-1   | 6.44| 7.18| 6.14
sorty-2   | 3.35| 3.63| 3.20
sorty-3   | 2.41| 2.53| 2.66
sorty-4   | 1.96| 1.92| 2.31
sortyLsw-1|13.92|15.45|14.51
sortyLsw-2| 7.27| 7.88| 7.58
sortyLsw-3| 5.13| 5.23| 6.07
sortyLsw-4| 4.05| 4.13| 5.37

Sorting random float32 array (in seconds):

Library(-Mxg)|R|X|i5
:---|---:|---:|---:
sort.Slice|17.49|18.18|17.21
sortutil  | 2.89| 3.07| 4.51
zermelo   | 4.64| 3.82| 3.17
sorty-1   | 7.78| 7.95| 6.96
sorty-2   | 4.00| 4.04| 3.64
sorty-3   | 2.78| 2.61| 2.84
sorty-4   | 2.39| 2.25| 2.54
sortyLsw-1|15.37|16.19|15.31
sortyLsw-2| 7.97| 8.25| 8.00
sortyLsw-3| 5.65| 5.55| 6.38
sortyLsw-4| 4.37| 4.34| 5.62

Sorting random string array (in seconds):

Library(-Mxg)|R|X|i5
:---|---:|---:|---:
sort.Slice| 8.28| 9.53| 8.46
sortutil  | 1.64| 2.11| 2.28
radix     | 4.27| 4.49| 3.50
sorty-1   | 5.82| 7.05| 6.32
sorty-2   | 2.89| 3.51| 3.40
sorty-3   | 2.00| 2.40| 3.07
sorty-4   | 1.64| 1.85| 2.84
sortyLsw-1| 7.50| 8.82| 8.10
sortyLsw-2| 3.75| 4.30| 4.38
sortyLsw-3| 2.58| 2.91| 3.99
sortyLsw-4| 2.19| 2.37| 3.75

### Testing & Parameter Tuning
First, make sure everything is fine:
```
go test -short -timeout 1h
```
You can tune `Mli,Hmli,Mlr` for your platform/cpu with (optimization flags):
```
go test -timeout 3h -gcflags '-B -wb=0 -smallframes' -ldflags '-s -w'
```
Now update `Mli,Hmli,Mlr` in sorty.go and compare your tuned sorty with others:
```
go test -short -timeout 1h -gcflags '-B -wb=0 -smallframes' -ldflags '-s -w'
```
Remember to build sorty (and your functions like `SortObjAsc()`) with the same
optimization flags you used for tuning.

### Support
If you use sorty and like it, please support via ETH:`0x464B840ee70bBe7962b90bD727Aac172Fa8B9C15`
