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

// allocate fill sort test
func afst(name string, srt func([]uint32)) []uint32 {
	ar := make([]uint32, N)
	fill(ar)

	now := time.Now()
	srt(ar)
	dur := time.Now().Sub(now)
	fmt.Println(name, "took", dur)

	if !IsSortedU4(ar) {
		tst.Fatal(name, "not sorted")
	}
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

	if !IsSortedF4(ar) {
		tst.Fatal(name, "not sorted")
	}
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

var name = []byte("sorty-00")

func Test1(t *testing.T) {
	tst = t
	fmt.Println("Sorting uint32")
	ar := afst("sort.Slice", func(ar []uint32) { sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] }) })

	ap := afst("sortutil", sortutil.Uint32s)
	compare(ap, ar)
	ap = afst("zermelo", func(ar []uint32) { zermelo.Sort(ar) })
	compare(ap, ar)

	for i := 2; i <= 16; i *= 2 {
		name[6] = byte(i/10) + '0'
		name[7] = byte(i%10) + '0'
		ap = afst(string(name), func(ar []uint32) { SortU4(ar, int32(i)) })
		compare(ap, ar)
	}
	fmt.Println()
}

func Test2(t *testing.T) {
	tst = t
	fmt.Println("Sorting float32")
	ar := afst2("sort.Slice", func(ar []float32) { sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] }) })

	ap := afst2("sortutil", sortutil.Float32s)
	compare2(ap, ar)
	ap = afst2("zermelo", func(ar []float32) { zermelo.Sort(ar) })
	compare2(ap, ar)

	for i := 2; i <= 16; i *= 2 {
		name[6] = byte(i/10) + '0'
		name[7] = byte(i%10) + '0'
		ap = afst2(string(name), func(ar []float32) { SortF4(ar, int32(i)) })
		compare2(ap, ar)
	}
	fmt.Println()
}
