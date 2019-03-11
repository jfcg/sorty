## sorty [![Go Report Card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty)
Type-specific concurrent sorting library

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation \(with [InsertionSort](https://en.wikipedia.org/wiki/Insertion_sort) as subroutine\) and does not require extra memory. Call corresponding Sort\*() to sort your slice in ascending order. For example:
```
sorty.SortS(your_string_slice, mx)
```
A Sort\*() function should not be called by multiple goroutines at the same time. mx is the maximum number of goroutines used for sorting simultaneously.

### 'go test' results on various computers
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts) and [zermelo](https://github.com/shawnsmithdev/zermelo).

Sorting uint32 array (in seconds):

Library|Acer TravelMate netbook with Intel Celeron N3160|Casper laptop with Intel Core i5-4210M|Server with AMD Ryzen 5 1600
:---|:---:|:---:|:---:
sort.Slice|144.40|68.48|67.58
sortutil  | 65.50|25.50|17.38
zermelo   | 39.40|11.35| 9.51
sorty-2   | 26.62|15.84|15.17
sorty-3   | 19.27|12.99|11.33
sorty-4   | 15.22|11.09| 9.20
sorty-5   | 15.16|11.15| 7.89

Sorting float32 array (in seconds):

Library|Acer TravelMate netbook with Intel Celeron N3160|Casper laptop with Intel Core i5-4210M|Server with AMD Ryzen 5 1600
:---|:---:|:---:|:---:
sort.Slice|147.80|69.50|71.73
sortutil  | 62.30|23.46|16.58
zermelo   | 31.20| 9.26|14.55
sorty-2   | 32.28|19.43|23.31
sorty-3   | 23.36|15.59|17.69
sorty-4   | 20.49|13.18|14.37
sorty-5   | 18.59|13.35|12.45

### Parameter Tuning
First, make sure everything is fine (prepend GOGC=30 to all if your ram <= 4 GiB):
```
go test -short -timeout 9h
```
You can tune Mli,Mlr for your platform/cpu with \(optimization flags\):
```
go test -timeout 99h -gcflags '-B -s' -ldflags '-s -w'
```
Now update Mli,Mlr in sorty.go and compare your tuned sorty with sortutil & zermelo:
```
go test -short -timeout 9h -gcflags '-B -s' -ldflags '-s -w'
```
Remember to build your sorty with the same flags you used for tuning.
