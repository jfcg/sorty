## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![go.dev ref](https://raw.githubusercontent.com/jfcg/.github/main/godev.svg)](https://pkg.go.dev/github.com/jfcg/sorty)
Type-specific, fast, efficient, concurrent/parallel sorting library.

sorty is a concurrent [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation (with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine). It is in-place and does not require extra memory. You can call corresponding `Sort*()` to rapidly sort your slices (in ascending order) or collections of objects. For example:
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
Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go version `1.15.7` on:

Machine|CPU|OS|Kernel
:---:|:---|:---|:---
R |Ryzen 1600   |Manjaro     |5.4.89
X |Xeon Haswell |Ubuntu 20.04|5.4.0-1034-gcp
i5|Core i5 4210M|Manjaro     |5.4.89

Sorting random uint32 array (in seconds):

Library(-Mxg)|R|X|i5
:---|---:|---:|---:
sort.Slice|16.05|17.93|16.10
sortutil  | 2.34| 2.95| 3.86
zermelo   | 2.00| 1.79| 1.09
sorty-1   | 6.50| 7.50| 6.14
sorty-2   | 3.37| 3.78| 3.20
sorty-3   | 2.48| 2.63| 2.65
sorty-4   | 1.88| 1.97| 2.29
sortyLsw-1|13.94|16.09|14.51
sortyLsw-2| 7.24| 8.15| 7.57
sortyLsw-3| 4.91| 5.43| 6.08
sortyLsw-4| 4.07| 4.28| 5.36

Sorting random float32 array (in seconds):

Library(-Mxg)|R|X|i5
:---|---:|---:|---:
sort.Slice|17.44|18.72|17.21
sortutil  | 2.86| 3.15| 4.50
zermelo   | 4.60| 3.98| 3.17
sorty-1   | 7.67| 8.26| 6.95
sorty-2   | 4.00| 4.20| 3.63
sorty-3   | 2.82| 2.73| 2.83
sorty-4   | 2.34| 2.31| 2.53
sortyLsw-1|15.34|16.81|15.30
sortyLsw-2| 7.96| 8.57| 8.01
sortyLsw-3| 5.65| 5.76| 6.44
sortyLsw-4| 4.31| 4.50| 5.61

Sorting random string array (in seconds):

Library(-Mxg)|R|X|i5
:---|---:|---:|---:
sort.Slice| 8.24| 9.82| 8.45
sortutil  | 1.64| 2.16| 2.28
radix     | 4.30| 4.71| 3.44
sorty-1   | 5.80| 7.37| 6.28
sorty-2   | 2.88| 3.67| 3.40
sorty-3   | 2.04| 2.53| 3.08
sorty-4   | 1.70| 1.94| 2.85
sortyLsw-1| 7.54| 9.16| 8.05
sortyLsw-2| 3.79| 4.55| 4.38
sortyLsw-3| 2.60| 3.02| 3.99
sortyLsw-4| 2.18| 2.45| 3.75

### Testing & Parameter Tuning
First, make sure everything is fine. Both [gc](https://golang.org) and [tinygo](https://tinygo.org) commands are provided:
```
tinygo test
go test -timeout 1h
```
You can tune `Mli,Hmli,Mlr` for your platform/cpu with (optimization flags):
```
tinygo test -opt 2 -no-debug -tags tuneparam
go test -timeout 3h -gcflags '-dwarf=0 -B -wb=0' -ldflags '-s -w' -tags tuneparam
```
Now update `Mli,Hmli,Mlr` in sorty.go and compare your tuned sorty with others:
```
tinygo test -opt 2 -no-debug
go test -timeout 1h -gcflags '-dwarf=0 -B -wb=0' -ldflags '-s -w'
```
Remember to build sorty (and your functions like `SortObjAsc()`) with the same
optimization flags you used for tuning.

### Support
If you use sorty and like it, please support via ETH:`0x464B840ee70bBe7962b90bD727Aac172Fa8B9C15`
