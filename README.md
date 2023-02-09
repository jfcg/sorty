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

sorty does not yet recognize partially sorted (sub-)slices to sort them faster (like pdqsort).

### Benchmarks
See `Green tick > QA / Tests > Details`. Testing and benchmarks are done with random inputs
via [jfcg/rng](https://github.com/jfcg/rng) library.

### Testing & Parameter Tuning
Run tests with:
```
go test -timeout 1h -v
```
You can tune `MaxLen*` for your platform/CPU with:
```
go test -timeout 3h -tags tuneparam
```
Now you can update `MaxLen*` in `maxc.go` and run tests again to see the improvements.
The parameters are already set to give good performance over different CPUs.
Also see `Green tick > QA / Tuning > Details`.

### Support
See [Contributing](./.github/CONTRIBUTING.md), [Security](./.github/SECURITY.md) and [Support](./.github/SUPPORT.md) guides. Also if you use sorty and like it, please support via [Github Sponsors](https://github.com/sponsors/jfcg) or:
- BTC:`bc1qr8m7n0w3xes6ckmau02s47a23e84umujej822e`
- ETH:`0x3a844321042D8f7c5BB2f7AB17e20273CA6277f6`
