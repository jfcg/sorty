package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import (
	"fmt"
	"github.com/shawnsmithdev/zermelo/zfloat32"
	"github.com/shawnsmithdev/zermelo/zuint32"
	"github.com/twotwotwo/sorts/sortutil"
	"math/rand"
	//"sort"
	"testing"
	"time"
	"unsafe"
)

const N = 1 << 28

var tst *testing.T
var name string

// fill sort test
func fst(sd int64, ar []uint32, srt func([]uint32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := N - 1; i >= 0; i-- {
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
func fst2(sd int64, ar []float32, srt func([]float32)) time.Duration {
	rn := rand.New(rand.NewSource(sd))
	for i := N - 1; i >= 0; i-- {
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

func compare(ar, ap []uint32) {
	if len(ap) <= 0 {
		return
	}
	if len(ar) != N || len(ap) != N {
		tst.Fatal(name, "length mismatch:", len(ap), len(ar))
	}

	for i := N - 1; i >= 0; i-- {
		if ap[i] != ar[i] {
			tst.Fatal(name, "values mismatch:", i, ap[i], ar[i])
		}
	}
}

// fst compare
func afc(srt func([]uint32), ar, ap []uint32) float64 {
	// take average time of two different sorts
	dur := fst(1, ar, srt)
	dur = (fst(2, ar, srt) + dur) / 2
	compare(ar, ap)

	sec := dur.Seconds()
	if testing.Short() {
		fmt.Printf("%s took %.2fs\n", name, sec)
	}
	return sec
}

func f2u(p *[]float32) []uint32 {
	return *(*[]uint32)(unsafe.Pointer(p))
}

// fst compare
func afc2(srt func([]float32), ar, ap []float32) float64 {
	// take average time of two different sorts
	dur := fst2(1, ar, srt)
	dur = (fst2(2, ar, srt) + dur) / 2
	compare(f2u(&ar), f2u(&ap))

	sec := dur.Seconds()
	if testing.Short() {
		fmt.Printf("%s took %.2fs\n", name, sec)
	}
	return sec
}

var srnm = []byte("sorty-0")

// return sum of Sort*() times for 2..5 goroutines
// compare with ap and among themselves
func sumt(ar, ap []uint32) float64 {
	sum := .0
	for i := 2; i < 6; i++ {
		srnm[6] = byte(i%10) + '0'
		name = string(srnm)
		sum += afc(func(ar []uint32) { SortU4(ar, uint32(i)) }, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return sum
}

// return sum of Sort*() times for 2..5 goroutines
// compare with ap and among themselves
func sumt2(ar, ap []float32) float64 {
	sum := .0
	for i := 2; i < 6; i++ {
		srnm[6] = byte(i%10) + '0'
		name = string(srnm)
		sum += afc2(func(ar []float32) { SortF4(ar, uint32(i)) }, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return sum
}

func TestShort(t *testing.T) {
	if !testing.Short() {
		t.SkipNow()
	}
	tst = t

	fmt.Println("Sorting uint32")
	//ar := make([]uint32, N) //too slow to test :P
	//name = "sort.Slice"
	//afc(func(ar []uint32) { sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] }) }, ar, nil)

	ar := make([]uint32, N)
	name = "sortutil"
	afc(sortutil.Uint32s, ar, nil)
	ap := make([]uint32, N)
	name = "zermelo"
	afc(zuint32.Sort, ap, ar)

	sumt(ap, ar)
	ar, ap = nil, nil

	fmt.Println("\nSorting float32")
	//aq := make([]float32, N)
	//name = "sort.Slice"
	//afc2(func(ar []float32) { sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] }) }, aq, nil)

	aq := make([]float32, N)
	name = "sortutil"
	afc2(sortutil.Float32s, aq, nil)
	as := make([]float32, N)
	name = "zermelo"
	afc2(zfloat32.Sort, as, aq)

	sumt2(as, aq)

	// SortI calls SortI4 (on 32-bit) or SortI8 (on 64-bit).
	SortI(iar, 3)
	if !IsSortedI(iar) {
		t.Fatal("SortI does not work")
	}
}

var iar = []int{
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	9, 8, 7, 6, 5, 4, 3, 2, 1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0,
	-9, -8, -7, -6, -5, -4, -3, -2, -1, 0, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0}

// optimize max array lengths for insertion sort/recursion (Mli,Mlr)
func TestOpt(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tst = t
	const cd = 3 // countdown iv

	fmt.Println("Sorting uint32")
	ar := make([]uint32, N)
	ap := make([]uint32, 0, N)

	mns, ic, rc := 9e9, cd, cd // min sum, insertion/recursion no-improvement countdowns
	for Mli = 16; ic > 0; Mli += 4 {

		for rc, Mlr = cd, 3*Mli-3; ; Mlr += 4 {
			sum := sumt(ar, ap)

			if sum < mns {
				mns, ic, rc = sum, cd+1, cd
				fmt.Printf("%d %d %.2fs\n", Mli, Mlr, mns)
			} else {
				rc--
				if rc <= 0 {
					break
				}
			}
		}
		ic--
	}
	ar, ap = nil, nil

	fmt.Println("\nSorting float32")
	aq := make([]float32, N)
	as := make([]float32, 0, N)

	mns, ic, rc = 9e9, cd, cd
	for Mli = 16; ic > 0; Mli += 4 {

		for rc, Mlr = cd, 3*Mli-3; ; Mlr += 4 {
			sum := sumt2(aq, as)

			if sum < mns {
				mns, ic, rc = sum, cd+1, cd
				fmt.Printf("%d %d %.2fs\n", Mli, Mlr, mns)
			} else {
				rc--
				if rc <= 0 {
					break
				}
			}
		}
		ic--
	}
}
