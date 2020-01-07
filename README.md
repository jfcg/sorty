## sorty [![Go Report Card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![GoDoc](https://godoc.org/github.com/jfcg/sorty?status.svg)](https://godoc.org/github.com/jfcg/sorty)
Type-specific concurrent / parallel sorting library

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
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix) with Go 1.13.5.

Sorting uint32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|15.99|17.37
sortutil  | 3.00| 3.49
zermelo   | 2.20| 1.85
sorty-2   | 3.20| 3.08
sorty-3   | 2.39| 2.24
sorty-4   | 1.98| 1.81
sorty-Col | 7.02| 6.98
sorty-Col2| 6.43| 6.33
sorty-lsw | 5.46| 5.43

Sorting float32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|17.57|17.96
sortutil  | 3.00| 4.13
zermelo   | 4.65| 3.40
sorty-2   | 4.03| 3.48
sorty-3   | 3.03| 2.52
sorty-4   | 2.49| 2.08
sorty-Col | 7.81| 7.27
sorty-Col2| 7.06| 6.52
sorty-lsw | 6.27| 5.63

Sorting string array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice| 8.54| 8.03
sortutil  | 1.97| 2.35
radix     | 4.63| 3.31
sorty-2   | 3.30| 3.47
sorty-3   | 2.53| 2.56
sorty-4   | 2.08| 2.12
sorty-Col2| 3.36| 3.44
sorty-lsw | 3.25| 3.39

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
