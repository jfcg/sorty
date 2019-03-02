## sorty [![Go Report Card](https://goreportcard.com/badge/github.com/jfcg/sorty)](https://goreportcard.com/report/github.com/jfcg/sorty)
Type-specific concurrent sorting library

Call corresponding Sort\*() to sort your slice in ascending order. For example:
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
sortutil  | 65.6|26.0|18.0
zermelo   | 38.6|11.7|10.3
sorty-02  | 31.3|18.5|18.7
sorty-04  | 18.4|13.1|12.3
sorty-08  | 17.9|12.7| 8.6
sorty-16  | 18.2|12.6| 5.8

Sorting float32 array (in seconds):

Library|Acer TravelMate netbook with Intel Celeron N3160|Casper laptop with Intel Core i5-4210M|Server with AMD Ryzen 5 1600
:---|:---:|:---:|:---:
sort.Slice|147.8|69.9|72.0
sortutil  | 60.8|24.0|16.6
zermelo   | 30.3| 8.7|14.5
sorty-02  | 36.8|21.5|24.1
sorty-04  | 28.2|14.3|14.8
sorty-08  | 22.9|14.2|10.8
sorty-16  | 25.3|14.1| 7.1
