## sorty [![Go Report Card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty) [![GoDoc](https://godoc.org/github.com/jfcg/sorty?status.svg)](https://godoc.org/github.com/jfcg/sorty)
Type-specific concurrent sorting library

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation \(with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine\) and does not require extra memory. Call corresponding Sort\*() to concurrently sort your slice (in ascending order) or collection. For example:
```
sorty.SortS(string_slice)
sorty.Sort(col)   // satisfies sort.Interface
sorty.Sort2(col2) // satisfies sorty.Collection2
```
Mxg (3 by default) is the maximum number of goroutines used for sorting per Sort\*() call.

### 'go test' results
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts), [zermelo](https://github.com/shawnsmithdev/zermelo) and [radix](https://github.com/yourbasic/radix).

Sorting uint32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|15.99|17.20
sortutil  | 3.01| 3.87
zermelo   | 2.18| 3.34
sorty-2   | 3.17| 3.08
sorty-3   | 2.45| 2.24
sorty-4   | 1.98| 1.83
sorty-Col | 7.05| 6.93
sorty-Col2| 6.40| 6.35

Sorting float32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|17.57|17.98
sortutil  | 3.13| 4.17
zermelo   | 4.65| 4.01
sorty-2   | 4.00| 3.47
sorty-3   | 3.03| 2.49
sorty-4   | 2.47| 2.06
sorty-Col | 7.76| 7.18
sorty-Col2| 7.01| 6.50

Sorting string array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice| 8.54| 8.66
sortutil  | 1.94| 2.65
radix     | 4.68| 4.42
sorty-2   | 3.30| 3.67
sorty-3   | 2.54| 2.73
sorty-4   | 2.06| 2.22

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
