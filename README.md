## sorty [![Go Report Card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![GoDoc](https://godoc.org/github.com/jfcg/sorty?status.svg)](https://godoc.org/github.com/jfcg/sorty)
Type-specific concurrent sorting library

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation \(with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine\) and does not require extra memory. Call corresponding Sort\*() to concurrently sort your slice (in ascending order) or collection. For example:
```
sorty.SortS(string_slice) // native slice
sorty.Sort(col)           // satisfies sort.Interface
sorty.Sort2(col2)         // satisfies sorty.Collection2
sorty.Sort3(n, lesswap)   // lesswap() function based
```
Mxg (3 by default) is the maximum number of goroutines used for sorting per Sort\*() call.
sorty uses [semantic](https://semver.org) versioning.

### 'go test' results
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix).

Sorting uint32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|15.99|17.37
sortutil  | 3.00| 3.52
zermelo   | 2.20| 1.85
sorty-2   | 3.19| 3.08
sorty-3   | 2.42| 2.24
sorty-4   | 2.00| 1.80
sorty-Col | 7.02| 6.96
sorty-Col2| 6.42| 6.33

Sorting float32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|17.57|17.96
sortutil  | 3.00| 4.15
zermelo   | 4.65| 3.40
sorty-2   | 4.02| 3.48
sorty-3   | 3.03| 2.50
sorty-4   | 2.49| 2.06
sorty-Col | 7.79| 7.23
sorty-Col2| 7.03| 6.52

Sorting string array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice| 8.54| 8.03
sortutil  | 1.96| 2.35
radix     | 4.64| 3.30
sorty-2   | 3.26| 3.47
sorty-3   | 2.54| 2.58
sorty-4   | 2.03| 2.12

### Parameter Tuning
First, make sure everything is fine:
```
go test -short -timeout 1h
```
You can tune Mli,Mlr for your platform/cpu with \(optimization flags\):
```
go test -timeout 2h -gcflags '-B -s' -ldflags '-s -w'
```
Now update Mli,Mlr in sorty.go and compare your tuned sorty with others:
```
go test -short -timeout 1h -gcflags '-B -s' -ldflags '-s -w'
```
Remember to build your sorty with the same flags you used for tuning.

### Support
If you use sorty and like it, please support via ETH:`0x464B840ee70bBe7962b90bD727Aac172Fa8B9C15`
