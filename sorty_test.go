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
func fst2(sd int64, ar []float32, srt func([]float32)) time.Duration {
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

func compare(ar, ap []uint32) {
	l := len(ap)
	if l <= 0 {
		return
	}
	if len(ar) != l {
		tst.Fatal(name, "length mismatch:", len(ar), l)
	}

	for i := l - 1; i >= 0; i-- {
		if ar[i] != ap[i] {
			tst.Fatal(name, "values mismatch:", i, ar[i], ap[i])
		}
	}
}

// average fst & compare
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

// average fst & compare
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
	s := .0
	for i := 2; i < 6; i++ {
		srnm[6] = byte(i + '0')
		name = string(srnm)
		s += afc(func(ar []uint32) { SortU4(ar, uint32(i)) }, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// return sum of Sort*() times for 2..5 goroutines
// compare with ap and among themselves
func sumt2(ar, ap []float32) float64 {
	s := .0
	for i := 2; i < 6; i++ {
		srnm[6] = byte(i + '0')
		name = string(srnm)
		s += afc2(func(ar []float32) { SortF4(ar, uint32(i)) }, ar, ap)
		ap, ar = ar, ap[:cap(ap)]
	}
	return s
}

// SortU4 and signal
func sas(sd int64, ar []uint32, ch chan bool) {
	fst(sd, ar, func(x []uint32) { SortU4(x, 2) })
	if ch != nil {
		ch <- false
	}
}

func TestShort(t *testing.T) {
	if !testing.Short() {
		t.SkipNow()
	}
	tst = t

	as := make([]float32, N)
	aq := make([]float32, N)
	ar, ap := f2u(&as), f2u(&aq)

	fmt.Println("Sorting uint32")
	//name = "sort.Slice" // takes too long
	//afc(func(ar []uint32) { sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] }) }, ar, nil)

	name = "sortutil"
	afc(sortutil.Uint32s, ar, nil)
	name = "zermelo"
	afc(zuint32.Sort, ap, ar)
	sumt(ap, ar) // sorty

	fmt.Println("\nSorting float32")
	//name = "sort.Slice"
	//afc2(func(ar []float32) { sort.Slice(ar, func(i, k int) bool { return ar[i] < ar[k] }) }, aq, nil)

	name = "sortutil"
	afc2(sortutil.Float32s, aq, nil)
	name = "zermelo"
	afc2(zfloat32.Sort, as, aq)
	sumt2(as, aq) // sorty

	// Is Sort*() multi-goroutine safe?
	fmt.Println("\nConcurrent calls to SortU4()")
	name = "multi"
	K, ch := N/2, make(chan bool, 1)

	go sas(3, ar[:K], ch)
	go sas(4, ar[K:], ch)
	go sas(3, ap[:K], ch)
	sas(4, ap[K:], nil)

	for i := 3; i > 0; i-- {
		<-ch // wait others
	}
	compare(ar[:K], ap[:K])
	compare(ar[K:], ap[K:])

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

// Mli,Mlr neighbors
var irN = []int{0, -4, 4, 0, -8, 0, 0, 8}

// Minimize fn() over Mli,Mlr grid
func findMin(fn func() float64) {
	Mli, Mlr = 24, 97

	// 3x3 grid of fn() values centered at Mli,Mlr
	var fv [9]float64
	fv[4] = fn() // center

	for {
		fmt.Printf("%d %d %.2fs\n", Mli, Mlr, fv[4])

		k, li, lr := 4, Mli, Mlr
		for i := 1; i < 8; i += 2 { // 4 neighbors

			if fv[i] > 0 { // known non-optimal?
				continue
			}

			Mli = li + irN[i/2] // peek neighbor
			Mlr = lr + irN[4+i/2]
			fv[i] = fn()

			if fv[i] < fv[k] { // better neighbor?
				k = i
			}
		}

		if k == 4 {
			break // center is best
		}

		Mli = li + irN[k/2] // switch to best neighbor
		Mlr = lr + irN[4+k/2]

		switch k { // update grid
		case 1: // up
			fv[6], fv[7], fv[8] = fv[3], fv[4], fv[5]
			fv[3], fv[4], fv[5] = fv[0], fv[1], fv[2]
			fv[0], fv[1], fv[2] = 0, 0, 0
		case 3: // left
			fv[2], fv[5], fv[8] = fv[1], fv[4], fv[7]
			fv[1], fv[4], fv[7] = fv[0], fv[3], fv[6]
			fv[0], fv[3], fv[6] = 0, 0, 0
		case 5: // right
			fv[0], fv[3], fv[6] = fv[1], fv[4], fv[7]
			fv[1], fv[4], fv[7] = fv[2], fv[5], fv[8]
			fv[2], fv[5], fv[8] = 0, 0, 0
		default: // down
			fv[0], fv[1], fv[2] = fv[3], fv[4], fv[5]
			fv[3], fv[4], fv[5] = fv[6], fv[7], fv[8]
			fv[6], fv[7], fv[8] = 0, 0, 0
		}
	}
}

// Optimize max array lengths for insertion sort/recursion (Mli,Mlr)
func TestOpt(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tst = t

	as := make([]float32, N)
	aq := make([]float32, 0, N)
	ar, ap := f2u(&as), f2u(&aq)

	fmt.Println("Sorting uint32")
	findMin(func() float64 { return sumt(ar, ap) })

	fmt.Println("\nSorting float32")
	findMin(func() float64 { return sumt2(as, aq) })
}
