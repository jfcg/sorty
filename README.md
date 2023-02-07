## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty/v2)](https://goreportcard.com/report/github.com/jfcg/sorty/v2) [![go.dev ref](https://pkg.go.dev/static/frontend/badge/badge.svg)](https://pkg.go.dev/github.com/jfcg/sorty/v2#pkg-overview)

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
[]uintptr, []float32, []float64, []string, [][]byte,
[]unsafe.Pointer, []*T // for any type T
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
version `1.17.8` on:

Machine|CPU|OS|Kernel
:---|:---|:---|:---
R6|Ryzen 1600   |Manjaro|5.10.105
i5|Core i5 4210M|Manjaro|5.10.105

Sorting uniformly distributed random uint32 slice (in seconds):

Library(-MaxGor)|R6|i5
:---|---:|---:
sort.Slice|12.06|14.01
  sortutil| 1.42| 3.12
   zermelo| 1.93| 1.12
   sorty-1| 6.18| 6.06
   sorty-2| 3.21| 3.18
   sorty-3| 2.17| 2.56
   sorty-4| 1.78| 2.26
sortyLsw-1|11.47|13.00
sortyLsw-2| 5.99| 6.80
sortyLsw-3| 4.08| 5.50
sortyLsw-4| 3.32| 4.78

Sorting normally distributed random float32 slice (in seconds):

Library(-MaxGor)|R6|i5
:---|---:|---:
sort.Slice|13.13|14.46
  sortutil| 1.99| 3.50
   zermelo| 4.51| 3.18
   sorty-1| 7.32| 6.86
   sorty-2| 3.89| 3.59
   sorty-3| 2.63| 2.78
   sorty-4| 2.29| 2.49
sortyLsw-1|12.83|13.60
sortyLsw-2| 6.76| 7.13
sortyLsw-3| 4.67| 5.63
sortyLsw-4| 3.88| 4.96

Sorting uniformly distributed random string slice (in seconds):

Library(-MaxGor)|R6|i5
:---|---:|---:
sort.Slice| 6.06| 7.05
  sortutil| 1.35| 1.94
   radix  | 4.26| 3.35
   sorty-1| 4.62| 5.30
   sorty-2| 2.41| 2.95
   sorty-3| 1.65| 2.73
   sorty-4| 1.50| 2.55
sortyLsw-1| 5.90| 6.77
sortyLsw-2| 3.12| 3.74
sortyLsw-3| 2.23| 3.37
sortyLsw-4| 1.98| 3.19

Sorting uniformly distributed random []byte slice (in seconds):

Library(-MaxGor)|R6|i5
:---|---:|---:
sort.Slice| 5.19| 6.20
   sorty-1| 3.32| 3.76
   sorty-2| 1.71| 2.05
   sorty-3| 1.25| 1.94
   sorty-4| 1.09| 1.80

Sorting uniformly distributed random string slice by length (in seconds):

Library(-MaxGor)|R6|i5
:---|---:|---:
sort.Slice| 2.99| 3.40
   sorty-1| 1.71| 1.91
   sorty-2| 0.95| 1.01
   sorty-3| 0.68| 0.86
   sorty-4| 0.57| 0.80

Sorting uniformly distributed random []byte slice by length (in seconds):

Library(-MaxGor)|R6|i5
:---|---:|---:
sort.Slice| 3.09| 3.47
   sorty-1| 1.18| 1.25
   sorty-2| 0.67| 0.67
   sorty-3| 0.47| 0.57
   sorty-4| 0.43| 0.54
</details>

### Testing & Parameter Tuning
<details><summary>Show testing & tuning</summary>

Run tests with:
```
go test -timeout 20m -v
```
You can tune `MaxLen*` for your platform/CPU with:
```
go test -timeout 2h -tags tuneparam
```
Now you can update `MaxLen*` in `maxc.go` and run tests again to see the improvements.
The parameters are already set to give good performance over different CPUs.
</details>

### Support
See [Contributing](./.github/CONTRIBUTING.md), [Security](./.github/SECURITY.md) and [Support](./.github/SUPPORT.md) guides. Also if you use sorty and like it, please support via [Github Sponsors](https://github.com/sponsors/jfcg) or:
- BTC:`bc1qr8m7n0w3xes6ckmau02s47a23e84umujej822e`
- ETH:`0x3a844321042D8f7c5BB2f7AB17e20273CA6277f6`
