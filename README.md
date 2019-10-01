## sorty [![Go Report Card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![GoDoc](https://godoc.org/github.com/jfcg/sorty?status.svg)](https://godoc.org/github.com/jfcg/sorty)
Type-specific concurrent sorting library

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation \(with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine\) and does not require extra memory. Call corresponding Sort\*() to concurrently sort your slice (in ascending order) or collection. For example:
```
sorty.SortS(string_slice)
sorty.Sort(col) // satisfies sort.Interface
```
Mxg (3 by default) is the maximum number of goroutines used for sorting per Sort\*() call.

### 'go test' results
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix).

Sorting uint32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|15.99|17.20
sortutil  | 2.94| 3.87
zermelo   | 2.20| 3.34
sorty-2   | 3.24| 3.08
sorty-3   | 2.41| 2.24
sorty-4   | 2.01| 1.83
sorty-Col | 6.98| 6.88

Sorting float32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|17.57|17.98
sortutil  | 3.12| 4.17
zermelo   | 4.64| 4.01
sorty-2   | 4.07| 3.44
sorty-3   | 3.04| 2.47
sorty-4   | 2.45| 2.05
sorty-Col | 7.63| 7.06

Sorting string array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice| 8.54| 8.66
sortutil  | 2.00| 2.63
radix     | 5.34| 4.45
sorty-2   | 3.28| 3.65
sorty-3   | 2.49| 2.73
sorty-4   | 2.01| 2.23

### Parameter Tuning
First, make sure everything is fine (prepend GOGC=30 to all if your ram <= 4 GiB):
```
go test -short -timeout 1h
```
You can tune Mli,Mlr for your platform/cpu with \(optimization flags\):
```
go test -timeout 1h -gcflags '-B -s' -ldflags '-s -w'
```
Now update Mli,Mlr in sorty.go and compare your tuned sorty with others:
```
go test -short -timeout 1h -gcflags '-B -s' -ldflags '-s -w'
```
Remember to build your sorty with the same flags you used for tuning.
