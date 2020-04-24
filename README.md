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
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go 1.14.2.

Sorting uint32 array (in seconds):

Library(-Mxg)|AMD Ryzen 5 1600
:---|:---:
sort.Slice|16.03
sortutil  | 3.00
zermelo   | 2.20
sorty-2   | 3.29
sorty-3   | 2.46
sorty-4   | 2.04
sortyLsw-2| 7.27
sortyLsw-3| 5.36
sortyLsw-4| 4.39

Sorting float32 array (in seconds):

Library(-Mxg)|AMD Ryzen 5 1600
:---|:---:
sort.Slice|17.43
sortutil  | 3.01
zermelo   | 4.69
sorty-2   | 4.07
sorty-3   | 3.00
sorty-4   | 2.43
sortyLsw-2| 8.06
sortyLsw-3| 5.94
sortyLsw-4| 4.80

Sorting string array (in seconds):

Library(-Mxg)|AMD Ryzen 5 1600
:---|:---:
sort.Slice| 8.72
sortutil  | 2.00
radix     | 4.83
sorty-2   | 3.24
sorty-3   | 2.48
sorty-4   | 2.07
sortyLsw-2| 4.10
sortyLsw-3| 3.12
sortyLsw-4| 2.60

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
