## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![go.dev ref](https://raw.githubusercontent.com/jfcg/.github/main/godev.svg)](https://pkg.go.dev/github.com/jfcg/sorty/v2)

sorty is a type-specific, fast, efficient, concurrent / parallel
[QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation (with an enhanced
[InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine).
It is in-place and does not require extra memory. You can call:
```
sorty.SortSlice(native_slice) // []int32, []float64 etc. in ascending order
sorty.SortLen(len_slice)      // []string or [][]T 'by length' in ascending order
sorty.Sort(n, lesswap)        // lesswap() based
```
If you have a pair of `Less()` and `Swap()`, then you can trivially write your
[`lesswap()`](https://pkg.go.dev/github.com/jfcg/sorty/v2#Sort) and sort your generic
collections using multiple CPU cores quickly.
sorty natively [sorts](https://pkg.go.dev/github.com/jfcg/sorty/v2#SortSlice)
```
[]int, []int32, []int64, []uint, []uint32, []uint64,
[]uintptr, []float32, []float64, []string, [][]byte
```
sorty also natively sorts `[]string` and `[][]T` (for any type `T`)
[by length](https://pkg.go.dev/github.com/jfcg/sorty/v2#SortLen).

sorty is stable (as in version), well-tested and pretty careful with resources & performance:
- `lesswap()` operates [**faster**](https://github.com/lynxkite/lynxkite/pull/141#issuecomment-779673635)
than [`sort.Interface`](https://pkg.go.dev/sort#Interface) on generic collections.
- For each `Sort*()` call, sorty uses up to [`MaxGor`](https://pkg.go.dev/github.com/jfcg/sorty/v2#pkg-variables)
(3 by default, including caller) concurrent goroutines and up to one channel.
- Goroutines and channel are created/used **only when necessary**.
- `MaxGor=1` (or a short input) yields single-goroutine sorting: no goroutines or channel will be created.
- `MaxGor` can be changed live, even during an ongoing `Sort*()` call.
- [`MaxLen*`](https://pkg.go.dev/github.com/jfcg/sorty/v2#pkg-constants) parameters are
tuned to get the best performance, see below.
- sorty API adheres to [semantic](https://semver.org) versioning.

### Benchmarks
Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts),
[zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go
version `1.17.1` on:

Machine|CPU|OS|Kernel
:---:|:---|:---|:---
R |Ryzen 1600    |Manjaro     |5.10.68
X |Xeon Broadwell|Ubuntu 20.04|5.11.0-1020-gcp
i5|Core i5 4210M |Manjaro     |5.10.68

Sorting uniformly distributed random uint32 slice (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice|12.18|16.78|13.98
  sortutil| 1.48| 3.42| 3.10
   zermelo| 2.10| 1.83| 1.12
   sorty-1| 5.83| 7.92| 6.06
   sorty-2| 3.07| 4.00| 3.17
   sorty-3| 2.30| 3.29| 2.62
   sorty-4| 1.82| 2.82| 2.27
sortyLsw-1|11.39|14.71|12.90
sortyLsw-2| 6.02| 7.49| 6.78
sortyLsw-3| 4.20| 5.96| 5.21
sortyLsw-4| 3.48| 5.18| 4.58

Sorting normally distributed random float32 slice (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice|13.28|17.37|14.43
  sortutil| 1.97| 3.95| 3.50
   zermelo| 4.50| 4.15| 3.18
   sorty-1| 7.21| 8.47| 6.84
   sorty-2| 3.78| 4.31| 3.57
   sorty-3| 2.60| 3.41| 2.78
   sorty-4| 2.27| 3.03| 2.50
sortyLsw-1|12.77|15.65|13.40
sortyLsw-2| 6.70| 8.00| 7.05
sortyLsw-3| 4.60| 6.20| 5.58
sortyLsw-4| 3.59| 5.50| 4.78

Sorting uniformly distributed random string slice (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice| 6.77| 8.96| 6.94
  sortutil| 1.31| 2.39| 1.99
   radix  | 4.90| 4.46| 3.41
   sorty-1| 4.80| 6.59| 5.09
   sorty-2| 2.46| 3.31| 2.81
   sorty-3| 1.80| 3.05| 2.63
   sorty-4| 1.42| 2.78| 2.51
sortyLsw-1| 6.28| 8.53| 6.73
sortyLsw-2| 3.24| 4.40| 3.66
sortyLsw-3| 2.20| 3.74| 3.29
sortyLsw-4| 1.95| 3.44| 3.12

Sorting uniformly distributed random []byte slice (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice| 6.87| 8.82| 7.04
   sorty-1| 4.07| 5.31| 4.28
   sorty-2| 2.12| 2.74| 2.38
   sorty-3| 1.49| 2.38| 2.15
   sorty-4| 1.30| 2.18| 1.97

Sorting uniformly distributed random string slice by length (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice| 3.37| 4.08| 3.45
   sorty-1| 1.61| 2.11| 1.67
   sorty-2| 0.87| 1.07| 0.88
   sorty-3| 0.61| 0.88| 0.72
   sorty-4| 0.53| 0.80| 0.67

Sorting uniformly distributed random []byte slice by length (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice| 3.55| 4.14| 3.51
   sorty-1| 1.18| 1.39| 1.12
   sorty-2| 0.63| 0.71| 0.60
   sorty-3| 0.45| 0.60| 0.51
   sorty-4| 0.39| 0.53| 0.46

### Testing & Parameter Tuning
First, make sure everything is fine:
```
go test -timeout 1h
```
You can tune `MaxLen*` for your platform/CPU with (optimization flags):
```
go test -timeout 4h -gcflags '-dwarf=0 -B -wb=0' -ldflags '-s -w' -tags tuneparam
```
Now update `MaxLen*` in `maxc.go`, uncomment imports & respective `mfc*()`
calls in `tmain_test.go` and compare your tuned sorty with other libraries:
```
go test -timeout 1h -gcflags '-dwarf=0 -B -wb=0' -ldflags '-s -w'
```
Remember to build sorty (and your functions like [`SortObjAsc()`](https://pkg.go.dev/github.com/jfcg/sorty/v2#Sort))
with the same optimization flags you used for tuning. `-B` flag is especially helpful.

### Support
If you use sorty and like it, please support via ETH:`0x464B840ee70bBe7962b90bD727Aac172Fa8B9C15`
