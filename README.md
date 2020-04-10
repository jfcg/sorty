## sorty [![go report card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![go.dev reference](/.github/godev.svg)](https://pkg.go.dev/github.com/jfcg/sorty?tab=doc)
Type-specific, fast concurrent / parallel sorting library.

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation (with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine) and does not require extra memory. Call corresponding Sort\*() to concurrently sort your slice (in ascending order) or collection. For example:
```
sorty.SortS(string_slice) // native slice
sorty.Sort(col)           // satisfies sort.Interface
sorty.Sort2(col2)         // satisfies sorty.Collection2
sorty.Sort3(n, lesswap)   // lesswap() function based
```
Mxg (3 by default) is the maximum number of goroutines used for sorting per Sort\*() call.
sorty uses [semantic](https://semver.org) versioning.

### Benchmarks
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go 1.14.2.

Sorting uint32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|16.03|17.37
sortutil  | 3.00| 3.49
zermelo   | 2.20| 1.85
sorty-2   | 3.29| 3.08
sorty-3   | 2.46| 2.24
sorty-4   | 2.04| 1.81
sorty-Col | 7.02| 6.98
sorty-Col2| 6.35| 6.33
sorty-lsw | 5.41| 5.43

Sorting float32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|17.43|17.96
sortutil  | 3.01| 4.13
zermelo   | 4.69| 3.40
sorty-2   | 4.07| 3.48
sorty-3   | 3.00| 2.52
sorty-4   | 2.43| 2.08
sorty-Col | 7.72| 7.27
sorty-Col2| 7.03| 6.52
sorty-lsw | 6.01| 5.63

Sorting string array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice| 8.72| 8.03
sortutil  | 2.00| 2.35
radix     | 4.83| 3.31
sorty-2   | 3.24| 3.47
sorty-3   | 2.48| 2.56
sorty-4   | 2.07| 2.12
sorty-Col2| 3.26| 3.44
sorty-lsw | 3.10| 3.39

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
