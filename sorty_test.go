package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import (
	"fmt"
	"github.com/twotwotwo/sorts/sortutil"
	"math/rand"
	"sort"
	"testing"
	"time"
)

const N = 1 << 28

var tst *testing.T

func fill() {
	rand.Seed(1)
	for i := len(ArU4) - 1; i >= 0; i-- {
		ArU4[i] = rand.Uint32()
	}
}

func isSorted(name string) {
	for i := len(ArU4) - 1; i > 0; i-- {
		if ArU4[i] < ArU4[i-1] {
			tst.Fatal(name, "not sorted:", i, ArU4[i], ArU4[i-1])
		}
	}
}

// allocate fill sort test
func afst(name string, srt func()) {
	ArU4 = make([]uint32, N)
	fill()

	now := time.Now()
	srt()
	dur := time.Now().Sub(now)
	fmt.Println(name, "took", dur)

	isSorted(name)
}

func compare(ar []uint32) {
	for i := len(ar) - 1; i >= 0; i-- {
		if ar[i] != ArU4[i] {
			tst.Fatal("Sorted arrays mismatch:", i, ar[i], ArU4[i])
		}
	}
}

func Test1(t *testing.T) {
	tst = t
	afst("sort.Slice", func() { sort.Slice(ArU4, func(i, k int) bool { return ArU4[i] < ArU4[k] }) })
	ar := ArU4

	afst("sorty", SortU4)
	compare(ar)
	afst("sortutil", func() { sortutil.Uint32s(ArU4) })
	compare(ar)
}
