## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![go.dev ref](https://raw.githubusercontent.com/jfcg/.github/main/godev.svg)](https://pkg.go.dev/github.com/jfcg/sorty)

sorty is a type-specific, fast, efficient, concurrent/parallel [QuickSort](https://en.wikipedia.org/wiki/Quicksort)
implementation (with an improved [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine).
It is in-place and does not require extra memory (other than efficient recursive calls and goroutines). You can call
corresponding `Sort*()` to rapidly sort your slices (in ascending order) or collections of objects. For example:
```
sorty.SortS(string_slice) // native slice
sorty.Sort(n, lesswap)    // lesswap() function based
```
If you have a pair of `Less()` and `Swap()`, then you can trivially write your
[`lesswap()`](https://pkg.go.dev/github.com/jfcg/sorty#Sort) and sort your generic
collections using multiple CPU cores quickly.

sorty is stable (as in version), well-tested and pretty careful with resources & performance:
- `lesswap()` operates [**faster**](https://github.com/lynxkite/lynxkite/pull/141#issuecomment-779673635)
than [`sort.Interface`](https://pkg.go.dev/sort#Interface) on generic collections.
- For each `Sort*()` call, sorty uses up to [`Mxg`](https://pkg.go.dev/github.com/jfcg/sorty#pkg-variables)
(3 by default, including caller) concurrent goroutines and up to one channel.
- Goroutines and channel are created/used **only when necessary**.
- `Mxg=1` (or a short input) yields single-goroutine sorting: No goroutines or channel will be created.
- `Mxg` can be changed live, even during an ongoing `Sort*()` call.
- `Mli,Hmli,Mlr` parameters are tuned to get the best performance, see below.
- sorty API adheres to [semantic](https://semver.org) versioning.

### Benchmarks
Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts),
[zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go
version `1.16.3` on:

Machine|CPU|OS|Kernel
:---:|:---|:---|:---
R |Ryzen 1600   |Manjaro     |5.4.108
X |Xeon Haswell |Ubuntu 20.04|5.4.0-1040-gcp
i5|Core i5 4210M|Manjaro     |5.4.108

Sorting uniformly distributed random uint32 array (in seconds):

Library(-Mxg)|R|X|i5
:---|---:|---:|---:
sort.Slice|15.82|18.44|16.02
  sortutil| 2.32| 3.08| 3.87
   zermelo| 2.00| 1.80| 1.09
   sorty-1| 6.35| 7.76| 6.13
   sorty-2| 3.32| 3.91| 3.19
   sorty-3| 2.15| 2.73| 2.65
   sorty-4| 1.89| 2.04| 2.27
sortyLsw-1|14.10|16.87|14.37
sortyLsw-2| 7.33| 8.56| 7.49
sortyLsw-3| 5.03| 5.78| 5.97
sortyLsw-4| 3.74| 4.44| 5.32

Sorting normally distributed random float32 array (in seconds):

Library(-Mxg)|R|X|i5
:---|---:|---:|---:
sort.Slice|17.44|19.53|17.15
  sortutil| 2.91| 3.11| 4.46
   zermelo| 4.50| 4.17| 3.18
   sorty-1| 7.59| 8.59| 6.99
   sorty-2| 4.00| 4.38| 3.66
   sorty-3| 2.71| 2.86| 2.83
   sorty-4| 2.36| 2.43| 2.53
sortyLsw-1|15.31|17.60|15.15
sortyLsw-2| 7.91| 8.96| 7.94
sortyLsw-3| 5.63| 6.00| 6.42
sortyLsw-4| 4.52| 4.71| 5.55

Sorting uniformly distributed random string array (in seconds):

Library(-Mxg)|R|X|i5
:---|---:|---:|---:
sort.Slice| 8.56|10.36| 8.42
  sortutil| 1.64| 2.39| 2.30
radix     | 4.29| 6.23| 3.38
   sorty-1| 5.76| 7.91| 6.23
   sorty-2| 2.89| 3.96| 3.41
   sorty-3| 1.97| 2.77| 3.10
   sorty-4| 1.76| 2.12| 2.86
sortyLsw-1| 7.52| 9.74| 8.08
sortyLsw-2| 3.76| 5.10| 4.37
sortyLsw-3| 2.68| 3.25| 3.96
sortyLsw-4| 2.17| 2.67| 3.77

### Testing & Parameter Tuning
First, make sure everything is fine:
```
go test -timeout 1h
```
You can tune `Mli,Hmli,Mlr` for your platform/CPU with (optimization flags):
```
go test -timeout 4h -gcflags '-dwarf=0 -B -wb=0' -ldflags '-s -w' -tags tuneparam
```
Now update `Mli,Hmli,Mlr` in `maxc.go`, uncomment imports & respective `mfc*()`
calls in `tmain_test.go` and compare your tuned sorty with other libraries:
```
go test -timeout 1h -gcflags '-dwarf=0 -B -wb=0' -ldflags '-s -w'
```
Remember to build sorty (and your functions like [`SortObjAsc()`](https://pkg.go.dev/github.com/jfcg/sorty#Sort))
with the same optimization flags you used for tuning. `-B` flag is especially helpful.

### Support
If you use sorty and like it, please support via ETH:`0x464B840ee70bBe7962b90bD727Aac172Fa8B9C15`
