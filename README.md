## sorty [![Go Report Card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty)
Type-specific concurrent sorting library

sorty is an in-place [QuickSort](https://en.wikipedia.org/wiki/Quicksort) implementation and does not require extra memory. Call corresponding Sort\*() to sort your slice in ascending order. For example:
```
sorty.SortS(your_string_slice, mx)
```
A Sort\*() function should not be called by multiple goroutines at the same time. mx is the maximum number of goroutines used for sorting simultaneously.

### 'go test' results on various computers
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts) and [zermelo](https://github.com/shawnsmithdev/zermelo).

Sorting uint32 array (in seconds):

Library|Acer TravelMate netbook with Intel Celeron N3160|Casper laptop with Intel Core i5-4210M|Server with AMD Ryzen 5 1600
:---|:---:|:---:|:---:
sort.Slice|144.4|68.6|67.2
sortutil  | 65.0|26.0|18.0
zermelo   | 38.0|12.0| 9.5
sorty-02  | 30.7|18.3|18.3
sorty-04  | 17.2|12.6|10.9
sorty-08  | 17.1|12.5| 8.2
sorty-16  | 17.1|12.4| 5.6

Sorting float32 array (in seconds):

Library|Acer TravelMate netbook with Intel Celeron N3160|Casper laptop with Intel Core i5-4210M|Server with AMD Ryzen 5 1600
:---|:---:|:---:|:---:
sort.Slice|147.8|69.9|72.0
sortutil  | 62.0|24.0|16.6
zermelo   | 30.3| 8.7|14.5
sorty-02  | 37.0|20.9|25.0
sorty-04  | 21.7|14.4|16.0
sorty-08  | 22.0|14.0|11.0
sorty-16  | 23.0|14.0| 7.8
