package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

const N = 1 << 27

var tst *testing.T

func fill() {
	rand.Seed(1)
	for i := len(Ar) - 1; i >= 0; i-- {
		Ar[i] = rand.Uint64()
	}
}

func isSorted(name string) {
	for i := len(Ar) - 1; i > 0; i-- {
		if Ar[i] < Ar[i-1] {
			tst.Fatal(name, "not sorted:", i, Ar[i], Ar[i-1])
		}
	}
}

// allocate fill sort test
func afst(name string, srt func()) {
	Ar = make([]uint64, N)
	fill()

	now := time.Now()
	srt()
	dur := time.Now().Sub(now)
	fmt.Println(name, "took", dur)

	isSorted(name)
}

func Test1(t *testing.T) {
	tst = t
	afst("sort.Slice", func() { sort.Slice(Ar, func(i, k int) bool { return Ar[i] < Ar[k] }) })
	ar := Ar
	afst("sorty", Sort)

	for i := N - 1; i >= 0; i-- {
		if ar[i] != Ar[i] {
			t.Fatal("Sorted arrays mismatch:", i, ar[i], Ar[i])
		}
	}
}
