package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import (
	"fmt"
	"github.com/shawnsmithdev/zermelo"
	"github.com/twotwotwo/sorts/sortutil"
	"math/rand"
	"sort"
	"testing"
	"time"
)

const N = 1 << 28

var tst *testing.T

func fill(ar []uint32) {
	rand.Seed(1)
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rand.Uint32()
	}
}

func fill2(ar []float32) {
	rand.Seed(2)
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = float32(rand.NormFloat64())
	}
}

func isSorted(name string, ar []uint32) {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			tst.Fatal(name, "not sorted:", i, ar[i], ar[i-1])
		}
	}
}

func isSorted2(name string, ar []float32) {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			tst.Fatal(name, "not sorted:", i, ar[i], ar[i-1])
		}
	}
}

// allocate fill sort test
func afst(name string, srt func([]uint32)) []uint32 {
	ar := make([]uint32, N)
	fill(ar)

	now := time.Now()
	srt(ar)
	dur := time.Now().Sub(now)
	fmt.Println(name, "took", dur)

	isSorted(name, ar)
	return ar
}

// allocate fill sort test
func afst2(name string, srt func([]float32)) []float32 {
	ar := make([]float32, N)
	fill2(ar)

	now := time.Now()
	srt(ar)
	dur := time.Now().Sub(now)
	fmt.Println(name, "took", dur)

	isSorted2(name, ar)
	return ar
}

func compare(ap, ar []uint32) {
	if len(ap) != len(ar) || len(ap) != N {
		tst.Fatal("Array length mismatch:", len(ap), len(ar))
	}

	for i := len(ar) - 1; i >= 0; i-- {
		if ap[i] != ar[i] {
			tst.Fatal("Array values mismatch:", i, ap[i], ar[i])
		}
	}
}

func compare2(ap, ar []float32) {
	if len(ap) != len(ar) || len(ap) != N {
		tst.Fatal("Array length mismatch:", len(ap), len(ar))
	}

	for i := len(ar) - 1; i >= 0; i-- {
		if ap[i] != ar[i] {
			tst.Fatal("Array values mismatch:", i, ap[i], ar[i])
		}
	}
}

func Test1(t *testing.T) {
	tst = t
	fmt.Println("Sorting uint32")
	ar := afst("sort.Slice", func(ar []uint32) { sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] }) })

	ap := afst("sortutil", sortutil.Uint32s)
	compare(ap, ar)
	ap = afst("zermelo", func(ar []uint32) { zermelo.Sort(ar) })
	compare(ap, ar)
	ap = afst("sorty", SortU4)
	compare(ap, ar)
}

func Test2(t *testing.T) {
	tst = t
	fmt.Println("Sorting float32")
	ar := afst2("sort.Slice", func(ar []float32) { sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] }) })

	ap := afst2("sortutil", sortutil.Float32s)
	compare2(ap, ar)
	ap = afst2("zermelo", func(ar []float32) { zermelo.Sort(ar) })
	compare2(ap, ar)
	ap = afst2("sorty", SortF4)
	compare2(ap, ar)
}
