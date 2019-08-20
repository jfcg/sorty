package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// Collection is same with standard library's sort.Interface. It represents a general
// collection of objects to be sorted. Note that Less() must be a strict ordering
// (like < or >), that is:
//  Less(i,k) && Less(k,r) => Less(i,r)
//  Less(i,k) => ! Less(k,i)
type Collection interface {
	// Len is the number of elements in the collection. First element is the 0-th one.
	Len() int
	// Less reports whether i-th element should sort before k-th element.
	Less(i, k int) bool
	// Swaps i-th and k-th elements.
	Swap(i, k int)
}

// IsSorted checks if ar is sorted.
func IsSorted(ar Collection) bool {
	for i := ar.Len() - 1; i > 0; i-- {
		if ar.Less(i, i-1) {
			return false
		}
	}
	return true
}

// insertion sort
func insertion(ar Collection, lo, hi int) {
	for h := lo + 1; h <= hi; h++ {
		for l := h; l > lo && ar.Less(l, l-1); l-- {
			ar.Swap(l, l-1)
		}
	}
}

// sort ar[l,m,h] and return swap status
func slmh(ar Collection, l, m, h int) int {
	// order l, h
	if ar.Less(h, l) {
		ar.Swap(h, l)
	}

	// order l, m, h
	if ar.Less(h, m) {
		ar.Swap(h, m)
		return 1
	}

	if ar.Less(m, l) {
		ar.Swap(m, l)
		return -1
	}
	return 0
}

// return pivot as median of five scattered values
func median(ar Collection, lo, hi int) int {
	mid := int(uint(lo+hi) >> 1)

	slmh(ar, lo, mid, hi)
	r := slmh(ar, lo+1, mid, hi-1)

	if r > 0 && ar.Less(mid, lo) {
		ar.Swap(mid, lo)
	}

	if r < 0 && ar.Less(hi, mid) {
		ar.Swap(hi, mid)
	}

	// here: ar[lo,lo+1] <= ar[mid] <= ar[hi-1,hi] as per Less()
	return mid
}

// Sort concurrently sorts ar.
func Sort(ar Collection) {
	arhi := ar.Len() - 1
	mli := Mli >> 2
	if arhi < mli {
		insertion(ar, 0, arhi)
		return
	}

	// number of sorting goroutines including this, end signal
	ng, done := uint32(1), make(chan bool, 1)
	var srt, gsrt func(int, int) // recursive & new-goroutine sort functions

	gsrt = func(lo, hi int) {
		srt(lo, hi)
		if atomic.AddUint32(&ng, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- false // we are the last, all done
		}
	}

	srt = func(lo, hi int) { // assumes hi-lo >= mli
		var l, h int
	start:
		l, h = lo+2, hi-2 // median handles lo,hi pairs

		// dl,dh (for avoiding unnecessary Less() calls) and pivot
		for dl, dh, pv := 1, -1, median(ar, lo, hi); l <= h; l, h = l+dl, h+dh {

			if dl == 0 {
				if ar.Less(h, pv) {
					ar.Swap(l, h)
					dl++
				}
				continue
			}

			if dh == 0 {
				if ar.Less(pv, l) {
					ar.Swap(l, h)
					dh--
				}
				continue
			}

			if ar.Less(h, pv) {
				if ar.Less(pv, l) {
					ar.Swap(l, h)
				} else {
					dh = 0
				}
			} else if ar.Less(pv, l) { // extend ranges in balance
				dl = 0
			}
		}

		if h-lo < hi-l {
			h, hi = hi, h // [lo,h] is the longer range
			l, lo = lo, l
		}

		// branches below are optimally laid out for fewer jumps
		// at least one short range?
		if hi-l < mli {
			insertion(ar, l, hi)

			if h-lo < mli { // two short ranges?
				insertion(ar, lo, h)
				return
			}
			hi = h
			goto start
		}

		// max goroutines? range not long enough for new goroutine?
		// not atomic but good enough
		if ng >= Mxg || hi-l < Mlr {
			srt(l, hi) // start a recursive sort on the shorter range
			hi = h
			goto start
		}

		if atomic.AddUint32(&ng, 1) == 0 { // increase goroutine counter
			panic("Sort: counter overflow")
		}
		go gsrt(lo, h) // start a new-goroutine sort on the longer range
		lo = l
		goto start
	}

	gsrt(0, arhi) // start master sort
	<-done
}
