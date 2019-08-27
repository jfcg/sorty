## sorty [![Go Report Card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty)
Type-specific concurrent sorting library

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation \(with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine\) and does not require extra memory. Call corresponding Sort\*() to concurrently sort your slice (in ascending order) or collection. For example:
```
sorty.SortS(string_slice)
sorty.Sort(col) // satisfies sort.Interface
```
Mxg (3 by default) is the maximum number of goroutines used for sorting per Sort\*() call.

### 'go test' results on various computers
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts) and [zermelo](https://github.com/shawnsmithdev/zermelo).

Sorting uint32 array (in seconds):

Library|Netbook with Intel Celeron N3160|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:|:---:
sort.Slice|34.75|15.97|17.22
sortutil  |10.18| 2.96| 3.86
zermelo   | 8.10| 2.21| 3.32
sorty-2   | 5.92| 3.64| 3.59
sorty-3   | 4.27| 2.76| 2.60
sorty-4   | 3.46| 2.28| 2.18
Sort(Col) |19.94| 7.59| 7.52

Sorting float32 array (in seconds):

Library|Netbook with Intel Celeron N3160|Server with AMD Ryzen 5 1600|Desktop with Intel Core i5-2400
:---|:---:|:---:|:---:
sort.Slice|36.49|17.43|17.98
sortutil  |11.62| 3.10| 4.19
zermelo   | 9.83| 4.65| 4.04
sorty-2   | 6.84| 4.45| 4.05
sorty-3   | 4.95| 3.30| 2.93
sorty-4   | 4.00| 2.79| 2.42
Sort(Col) |19.72| 8.21| 7.66

### Parameter Tuning
First, make sure everything is fine (prepend GOGC=30 to all if your ram <= 4 GiB):
```
go test -short -timeout 1h
```
You can tune Mli,Mlr for your platform/cpu with \(optimization flags\):
```
go test -timeout 1h -gcflags '-B -s' -ldflags '-s -w'
```
Now update Mli,Mlr in sorty.go and compare your tuned sorty with sortutil & zermelo:
```
go test -short -timeout 1h -gcflags '-B -s' -ldflags '-s -w'
```
Remember to build your sorty with the same flags you used for tuning.
