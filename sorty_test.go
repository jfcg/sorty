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

// fill sort test
func fst(name string, sd int64, ar []uint32, srt func([]uint32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = rn.Uint32()
	}

	now := time.Now()
	srt(ar)
	dur := time.Now().Sub(now)

	if !IsSortedU4(ar) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

// fill sort test
func fst2(name string, sd int64, ar []float32, srt func([]float32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := len(ar) - 1; i >= 0; i-- {
		ar[i] = float32(rn.NormFloat64())
	}

	now := time.Now()
	srt(ar)
	dur := time.Now().Sub(now)

	if !IsSortedF4(ar) {
		tst.Fatal(name, "not sorted")
	}
	return dur
}

// allocate fst compare
func afc(name string, srt func([]uint32), ap []uint32) ([]uint32, float64) {
	ar := make([]uint32, N)
	dur := fst(name, 1, ar, srt) // take average time of two different sorts
	dur = (fst(name, 2, ar, srt) + dur) / 2

	sec := dur.Seconds()
	if testing.Short() {
		fmt.Printf("%s took %.2fs\n", name, sec)
	}

	if ap != nil { // compare
		if len(ap) != len(ar) || len(ap) != N {
			tst.Fatal(name, "length mismatch:", len(ap), len(ar))
		}

		for i := len(ar) - 1; i >= 0; i-- {
			if ap[i] != ar[i] {
				tst.Fatal(name, "values mismatch:", i, ap[i], ar[i])
			}
		}
	}
	return ar, sec
}

// allocate fst compare
func afc2(name string, srt func([]float32), ap []float32) ([]float32, float64) {
	ar := make([]float32, N)
	dur := fst2(name, 1, ar, srt) // take average time of two different sorts
	dur = (fst2(name, 2, ar, srt) + dur) / 2

	sec := dur.Seconds()
	if testing.Short() {
		fmt.Printf("%s took %.2fs\n", name, sec)
	}

	if ap != nil { // compare
		if len(ap) != len(ar) || len(ap) != N {
			tst.Fatal(name, "length mismatch:", len(ap), len(ar))
		}

		for i := len(ar) - 1; i >= 0; i-- {
			if ap[i] != ar[i] {
				tst.Fatal(name, "values mismatch:", i, ap[i], ar[i])
			}
		}
	}
	return ar, sec
}

var name = []byte("sorty-0")

// return sum of Sort*() times for 2..5 goroutines
// compare with ap and among themselves
func sumt(ap []uint32) float64 {
	sum, t := .0, .0
	for i := 2; i < 6; i++ {
		name[6] = byte(i%10) + '0'
		ap, t = afc(string(name), func(ar []uint32) { SortU4(ar, uint32(i)) }, ap)
		sum += t
	}
	return sum
}

// return sum of Sort*() times for 2..5 goroutines
// compare with ap and among themselves
func sumt2(ap []float32) float64 {
	sum, t := .0, .0
	for i := 2; i < 6; i++ {
		name[6] = byte(i%10) + '0'
		ap, t = afc2(string(name), func(ar []float32) { SortF4(ar, uint32(i)) }, ap)
		sum += t
	}
	return sum
}

func Test1(t *testing.T) {
	if !testing.Short() {
		t.SkipNow()
	}
	tst = t
	fmt.Println("Sorting uint32")

	ar, _ := afc("sort.Slice", func(ar []uint32) {
		sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] })
	}, nil)
	_, _ = afc("sortutil", sortutil.Uint32s, ar)
	_, _ = afc("zermelo", func(ar []uint32) { zermelo.Sort(ar) }, ar)

	sumt(ar)
	fmt.Println()
}

func Test2(t *testing.T) {
	if !testing.Short() {
		t.SkipNow()
	}
	tst = t
	fmt.Println("Sorting float32")

	ar, _ := afc2("sort.Slice", func(ar []float32) {
		sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] })
	}, nil)
	_, _ = afc2("sortutil", sortutil.Float32s, ar)
	_, _ = afc2("zermelo", func(ar []float32) { zermelo.Sort(ar) }, ar)

	sumt2(ar)
	fmt.Println()
}

var iar = []int{
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0}

// SortI calls SortI4 (on 32-bit) or SortI8 (on 64-bit).
func TestI(t *testing.T) {
	SortI(iar, 3)
	if !IsSortedI(iar) {
		t.Fatal("SortI does not work")
	}
}

// optimize max array length for insertion sort (Mli)
func TestOpt(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tst = t
	fmt.Println("Sorting uint32")

	mns := 9e9
	for Mli = 10; Mli < 42; Mli += 2 {

		sum := sumt(nil)
		if sum < mns {
			mns = sum
			fmt.Printf("%d %.2fs\n", Mli, mns)
		}
	}
	fmt.Println()
}

// optimize max array length for insertion sort (Mli)
func TestOpt2(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tst = t
	fmt.Println("Sorting float32")

	mns := 9e9
	for Mli = 10; Mli < 42; Mli += 2 {

		sum := sumt2(nil)
		if sum < mns {
			mns = sum
			fmt.Printf("%d %.2fs\n", Mli, mns)
		}
	}
	fmt.Println()
}
