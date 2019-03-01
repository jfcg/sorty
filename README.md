## sorty
Type-specific concurrent sorting library

Call corresponding Sort\*() to sort your slice in ascending order. For example:
```
sorty.SortS(your_string_slice, mx)
```
A Sort\*() function should not be called by multiple goroutines at the same time. mx is the maximum number of goroutines used for sorting simultaneously.

### 'go test' results on various computers
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts) and [zermelo](https://github.com/shawnsmithdev/zermelo).

Acer TravelMate netbook with Intel Celeron N3160:
```
Sorting uint32
sort.Slice took 145s
sortutil took 65s
zermelo took 39s
sorty-02 took 31s
sorty-04 took 18.7s
sorty-08 took 17.4s
sorty-16 took 17.7s

Sorting float32
sort.Slice took 148s
sortutil took 60.5s
zermelo took 31.7s
sorty-02 took 37.5s
sorty-04 took 21.9s
sorty-08 took 21.1s
sorty-16 took 21.4s
```

Casper laptop with Intel Core i5-4210M:
```
Sorting uint32
sort.Slice took 68s
sortutil took 26s
zermelo took 12s
sorty-02 took 19s
sorty-04 took 13.5s
sorty-08 took 13s
sorty-16 took 12.9s

Sorting float32
sort.Slice took 69s
sortutil took 24s
zermelo took 8.6s
sorty-02 took 22.2s
sorty-04 took 15s
sorty-08 took 14.5s
sorty-16 took 14.5s
```

Server with AMD Ryzen 5 1600:
```
Sorting uint32
sort.Slice took 68s
sortutil took 17.6s
zermelo took 9.1s
sorty-02 took 18.2s
sorty-04 took 12s
sorty-08 took 8.6s
sorty-16 took 6.3s

Sorting float32
sort.Slice took 71.7s
sortutil took 16.5s
zermelo took 14.6s
sorty-02 took 25.3s
sorty-04 took 15.9s
sorty-08 took 11.3s
sorty-16 took 7.7s
```
