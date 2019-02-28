## sorty
Type-specific concurrent sorting library

Call corresponding Sort\*() to sort your slice in ascending order. For example:
```
sorty.SortS(your_string_slice)
```
A Sort\*() function should not be called by multiple goroutines at the same time. There is no limit on the number of goroutines to be created \(could be many thousands depending on data\), though sorty does it sparingly.

### 'go test' results on various computers
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort), [sortutil](https://github.com/twotwotwo/sorts) and [zermelo](https://github.com/shawnsmithdev/zermelo).

Acer TravelMate netbook with Intel Celeron N3160:
```
Sorting uint32
sort.Slice took 145s
sortutil took 65s
zermelo took 39s
sorty took 19.4s

Sorting float32
sort.Slice took 148s
sortutil took 60.5s
zermelo took 31.7s
sorty took 22.5s
```

Casper laptop with Intel Core i5-4210M:
```
Sorting uint32
sort.Slice took 68s
sortutil took 26.2s
zermelo took 12s
sorty took 13.5s

Sorting float32
sort.Slice took 69s
sortutil took 25s
zermelo took 8.6s
sorty took 14.6s
```

Server with AMD Ryzen 5 1600:
```
Sorting uint32
sort.Slice took 68s
sortutil took 17.6s
zermelo took 9.1s
sorty took 5.2s

Sorting float32
sort.Slice took 71.7s
sortutil took 16.5s
zermelo took 14.6s
sorty took 6.2s
```
