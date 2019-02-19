## sorty
Type-specific concurrent sorting library

Call corresponding Sort\*() to sort your slice in ascending order. For example:
```
sorty.SortS(your_string_slice)
```
A Sort\*() function should not be called by multiple goroutines at the same time. There is no limit on the number of goroutines to be created (could be many thousands depending on data), though sorty does it sparingly.

### 'go test' results on various computers
All computers run 64-bit Manjaro Linux. Comparing against [sort.Slice](https://golang.org/pkg/sort) and [sortutil](https://github.com/twotwotwo/sorts).

Acer TravelMate netbook with Intel Celeron N3160:
```
sort.Slice took 145s
sortutil took 65s
sorty took 19.4s
```

Casper laptop with Intel Core i5-4210M:
```
sort.Slice took 68s
sortutil took 26.2s
sorty took 14.2s
```

Server with AMD Ryzen 5 1600:
```
sort.Slice took 68s
sortutil took 17.6s
sorty took 5.2s
```
