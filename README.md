## sorty
Type-specific concurrent sorting library

Assign your slice to Ar\* and call Sort\*() to sort ascending. For example:
```
sorty.ArS = your_string_slice
sorty.SortS()
```
There is no limit on number of goroutines to create, though sorty does it sparingly.

### 'go test' results on various computers (64-bit Manjaro Linux)

Acer TravelMate netbook with Intel Celeron N3160:
```
[sort.Slice](https://golang.org/pkg/sort) took 145s
sorty took 19.4s
[sortutil](https://github.com/twotwotwo/sorts) took 65s
```

Casper laptop with Intel Core i5-4210M:
```
sort.Slice took 68s
sorty took 14.2s
sortutil took 26.2s
```

Server with AMD Ryzen 5 1600:
```
sort.Slice took 68s
sorty took 5.2s
sortutil took 17.6s
```
