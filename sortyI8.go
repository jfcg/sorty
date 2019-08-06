package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// IsSortedI8 checks if ar is sorted in ascending order.
func IsSortedI8(ar []int64) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

// insertion sort
func insertionI8(ar []int64) {
	for h := 1; h < len(ar); h++ {
		v, l := ar[h], h-1
		if v < ar[l] {
			for {
				ar[l+1] = ar[l]
				l--
				if l < 0 || v >= ar[l] {
					break
				}
			}
			ar[l+1] = v
		}
	}
}

// given vl <= vh, inserts pv in the middle
// returns vl <= pv <= vh
func ipI8(pv, vl, vh int64) (a, b, c int64, r int) {
	if pv > vh {
		vh, pv = pv, vh
		r = 1
	} else if pv < vl {
		vl, pv = pv, vl
		r = -1
	}
	return vl, pv, vh, r
}

// return pivot as median of five scattered values
func medianI8(ar []int64) int64 {
	// lo, mid, hi
	h := len(ar) - 1
	m := h >> 1
	vl, pv, vh := ar[0], ar[m], ar[h]

	// intermediates
	a, b := m>>1, int(uint(m+h)>>1) // avoid overflow
	va, vb := ar[a], ar[b]

	// put lo, mid, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipI8(pv, vl, vh)

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipI8(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipI8(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipI8(vh, pv, vb)
	}

	// here: vl, va <= pv <= vb, vh
	ar[0], ar[m], ar[h] = vl, pv, vh
	ar[a], ar[b] = va, vb
	return pv
}

// SortI8 concurrently sorts ar in ascending order. mx is the maximum number
// of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortI8(ar []int64, mx uint32) {
	if len(ar) <= Mli {
		insertionI8(ar)
		return
	}

	ng, mx := uint32(1), sat(mx) // number of sorting goroutines including this, max limit
	done := make(chan bool, 1)   // end signal
	var srt, gsrt func(int, int) // recursive & new-goroutine sort functions

	gsrt = func(lo, hi int) {
		srt(lo, hi)

		if atomic.AddUint32(&ng, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- false // we are the last, all done
		}
	}

	srt = func(lo, hi int) { // assumes hi-lo >= Mli
		var l, h int
	start:
		l, h = lo+1, hi-1 // medianI8 handles lo,hi positions

		for pv := medianI8(ar[lo : hi+1]); l <= h; {
			swap := true
			if ar[h] >= pv { // extend ranges in balance
				h--
				swap = false
			}
			if ar[l] <= pv {
				l++
				swap = false
			}

			if swap {
				ar[l], ar[h] = ar[h], ar[l]
				h--
				l++
			}
		}

		if h-lo < hi-l {
			h, hi = hi, h // [lo,h] is the bigger range
			l, lo = lo, l
		}

		if hi-l >= Mli { // two big ranges?

			// max goroutines? range not big enough for new goroutine?
			// not atomic but good enough
			if ng >= mx || hi-l < Mlr {
				srt(l, hi) // start a recursive sort on the smaller range
				hi = h
				goto start
			}

			if atomic.AddUint32(&ng, 1) == 0 { // increase goroutine counter
				panic("SortI8: counter overflow")
			}
			go gsrt(lo, h) // start a new goroutine on the bigger range
			lo = l
			goto start
		}

		insertionI8(ar[l : hi+1])

		if h-lo < Mli { // two small ranges?
			insertionI8(ar[lo : h+1])
			return
		}

		hi = h
		goto start
	}

	gsrt(0, len(ar)-1) // start sort
	<-done
}
