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
sortutil  | 2.97| 3.87
zermelo   | 2.20| 3.36
sorty-2   | 3.31| 3.12
sorty-3   | 2.47| 2.25
sorty-4   | 2.00| 1.84
sorty-Col | 7.05| 6.91

Sorting float32 array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:
sort.Slice|17.57|17.98
sortutil  | 3.12| 4.17
zermelo   | 4.64| 4.00
sorty-2   | 4.05| 3.43
sorty-3   | 3.01| 2.46
sorty-4   | 2.44| 2.04
sorty-Col | 7.72| 7.06

Sorting string array (in seconds):

Library|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:|:---:
sort.Slice| 8.54| 8.66s
sortutil  | 2.01| 2.63s
radix     | 4.60| 4.45s
sorty-2   | 3.27| 3.65s
sorty-3   | 2.48| 2.73s
sorty-4   | 2.02| 2.23s

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
