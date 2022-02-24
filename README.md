## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty/v2)](https://goreportcard.com/report/github.com/jfcg/sorty/v2) [![go.dev ref](https://raw.githubusercontent.com/jfcg/.github/main/godev.svg)](https://pkg.go.dev/github.com/jfcg/sorty/v2#pkg-overview)

sorty is a type-specific, fast, efficient, concurrent / parallel sorting
library. It is an innovative [QuickSort](https://en.wikipedia.org/wiki/Quicksort)
implementation, hence in-place and does not require extra memory. You can call:
```go
import "github.com/jfcg/sorty/v2"

sorty.SortSlice(native_slice) // []int, []float64, []string etc. in ascending order
sorty.SortLen(len_slice)      // []string or [][]T 'by length' in ascending order
sorty.Sort(n, lesswap)        // lesswap() based
```
If you have a pair of `Less()` and `Swap()`, then you can trivially write your
[`lesswap()`](https://pkg.go.dev/github.com/jfcg/sorty/v2#Sort) and sort your generic
collections using multiple CPU cores quickly.

sorty natively [sorts](https://pkg.go.dev/github.com/jfcg/sorty/v2#SortSlice) any type equivalent to
```go
[]int, []int32, []int64, []uint, []uint32, []uint64,
[]uintptr, []float32, []float64, []string, [][]byte
```
sorty also natively sorts any type equivalent to `[]string` or `[][]T` (for any type `T`)
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
<details><summary>Show benchmarks</summary>

Comparing against [sort.Slice](https://pkg.go.dev/sort#Slice), [sortutil](https://github.com/twotwotwo/sorts),
[zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go
version `1.17.3` on:

Machine|CPU|OS|Kernel
:---:|:---|:---|:---
R |Ryzen 1600    |Manjaro     |5.10.70
X |Xeon Broadwell|Ubuntu 20.04|5.11.0-1022-gcp
i5|Core i5 4210M |Manjaro     |5.10.70

Sorting uniformly distributed random uint32 slice (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice|12.09|16.59|13.98
  sortutil| 1.48| 3.42| 3.10
   zermelo| 2.10| 1.83| 1.12
   sorty-1| 5.96| 7.77| 6.03
   sorty-2| 3.18| 3.98| 3.17
   sorty-3| 2.28| 3.18| 2.56
   sorty-4| 1.86| 2.81| 2.29
sortyLsw-1|11.30|15.15|13.07
sortyLsw-2| 5.86| 7.69| 6.83
sortyLsw-3| 4.08| 6.14| 5.34
sortyLsw-4| 3.36| 5.31| 4.60

Sorting normally distributed random float32 slice (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice|13.19|17.16|14.44
  sortutil| 1.97| 3.95| 3.50
   zermelo| 4.50| 4.15| 3.18
   sorty-1| 7.21| 8.48| 6.84
   sorty-2| 3.76| 4.32| 3.58
   sorty-3| 2.48| 3.35| 2.78
   sorty-4| 2.12| 3.04| 2.47
sortyLsw-1|12.58|15.51|13.60
sortyLsw-2| 6.61| 7.90| 7.11
sortyLsw-3| 4.78| 6.44| 5.65
sortyLsw-4| 3.73| 5.56| 4.78

Sorting uniformly distributed random string slice (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice| 6.20| 8.87| 6.95
  sortutil| 1.31| 2.39| 1.99
   radix  | 4.90| 4.46| 3.41
   sorty-1| 4.69| 6.60| 5.15
   sorty-2| 2.41| 3.35| 2.86
   sorty-3| 1.75| 3.00| 2.69
   sorty-4| 1.47| 2.80| 2.51
sortyLsw-1| 5.72| 8.47| 6.63
sortyLsw-2| 2.98| 4.32| 3.67
sortyLsw-3| 2.18| 3.78| 3.28
sortyLsw-4| 1.79| 3.44| 3.08

Sorting uniformly distributed random []byte slice (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice| 6.00| 8.81| 7.20
   sorty-1| 3.25| 4.59| 3.63
   sorty-2| 1.68| 2.31| 1.99
   sorty-3| 1.23| 2.11| 1.90
   sorty-4| 1.04| 1.92| 1.76

Sorting uniformly distributed random string slice by length (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice| 2.96| 4.04| 3.40
   sorty-1| 1.55| 2.16| 1.72
   sorty-2| 0.87| 1.11| 0.92
   sorty-3| 0.63| 0.94| 0.78
   sorty-4| 0.51| 0.83| 0.69

Sorting uniformly distributed random []byte slice by length (in seconds):

Library(-MaxGor)|R|X|i5
:---|---:|---:|---:
sort.Slice| 2.98| 4.08| 3.56
   sorty-1| 1.09| 1.39| 1.11
   sorty-2| 0.60| 0.72| 0.60
   sorty-3| 0.42| 0.58| 0.52
   sorty-4| 0.39| 0.54| 0.47
</details>

### Testing & Parameter Tuning
First, make sure everything is fine:
```sh
go test -timeout 1h
```
You can tune `MaxLen*` for your platform/CPU with (optimization flags):
```sh
go test -timeout 4h -gcflags '-dwarf=0 -B' -ldflags '-s -w' -tags tuneparam
```
Now update `MaxLen*` in `maxc.go`, uncomment imports & respective `mfc*()`
calls in `tmain_test.go` and compare your tuned sorty with other libraries:
```sh
go test -timeout 1h -gcflags '-dwarf=0 -B' -ldflags '-s -w'
```
Remember to build sorty (and your functions like [`SortObjAsc()`](https://pkg.go.dev/github.com/jfcg/sorty/v2#Sort))
with the same optimization flags you used for tuning. `-B` flag is especially helpful.

### Support
If you use sorty and like it, please support via:
- BTC:`bc1qr8m7n0w3xes6ckmau02s47a23e84umujej822e`
- ETH:`0x3a844321042D8f7c5BB2f7AB17e20273CA6277f6`
