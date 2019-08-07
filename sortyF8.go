package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// IsSortedF8 checks if ar is sorted in ascending order.
func IsSortedF8(ar []float64) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

// insertion sort
func insertionF8(ar []float64) {
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
func ipF8(pv, vl, vh float64) (a, b, c float64, r int) {
	if pv > vh {
		return vl, vh, pv, 1
	} else if pv < vl {
		return pv, vl, vh, -1
	}
	return vl, pv, vh, 0
}

// return pivot as median of five scattered values
func medianF8(ar []float64) float64 {
	// lo, mid, hi
	h := len(ar) - 1
	m := h >> 1
	vl, pv, vh := ar[0], ar[m], ar[h]

	// put lo, mid, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipF8(pv, vl, vh)

	// intermediates
	a, b := m>>1, int(uint(m+h)>>1) // avoid overflow
	va, vb := ar[a], ar[b]

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipF8(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipF8(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipF8(vh, pv, vb)
	}

	// here: vl, va <= pv <= vb, vh
	ar[a], ar[m], ar[b] = va, pv, vb
	ar[0], ar[h] = vl, vh // update lo,hi positions last for better locality
	return pv
}

// SortF8 concurrently sorts ar in ascending order.
func SortF8(ar []float64) {
	if len(ar) <= Mli {
		insertionF8(ar)
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

	srt = func(lo, hi int) { // assumes hi-lo >= Mli
		var l, h int
	start:
		l, h = lo+1, hi-1 // medianF8 handles lo,hi positions

		for pv := medianF8(ar[lo : hi+1]); l <= h; {
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
			h, hi = hi, h // [lo,h] is the longer range
			l, lo = lo, l
		}

		// branches below are optimally laid out for less # of jumps
		// at least one short range?
		if hi-l < Mli {
			insertionF8(ar[l : hi+1])

			if h-lo < Mli { // two short ranges?
				insertionF8(ar[lo : h+1])
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
			panic("SortF8: counter overflow")
		}
		go gsrt(lo, h) // start a new-goroutine sort on the longer range
		lo = l
		goto start
	}

	gsrt(0, len(ar)-1) // start master sort
	<-done
}
