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

// sort and return vl,pv,vh & swap status
func slmhF8(vl, pv, vh float64) (a, b, c float64, r int) {
	// order vl, vh
	if vh < vl {
		vh, vl = vl, vh
	}

	// order vl, pv, vh
	if vh < pv {
		return vl, vh, pv, 1
	}

	if pv < vl {
		return pv, vl, vh, -1
	}
	return vl, pv, vh, 0
}

// return pivot as median of five scattered values
func medianF8(ar []float64) float64 {
	// lo, mid, hi
	h := len(ar) - 1
	m := h >> 1
	vl, va, pv, vb, vh := ar[0], ar[1], ar[m], ar[h-1], ar[h]

	vl, pv, vh, _ = slmhF8(vl, pv, vh)
	va, pv, vb, r := slmhF8(va, pv, vb)

	// if pivot was out of [va, vb]
	if r > 0 && pv < vl {
		pv, vl = vl, pv
	}

	if r < 0 && vh < pv {
		vh, pv = pv, vh
	}

	// here: vl, va <= pv <= vb, vh
	ar[0], ar[1], ar[m], ar[h-1], ar[h] = vl, va, pv, vb, vh
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
	start:
		l, h, pv := lo+2, hi-2, medianF8(ar[lo:hi+1]) // medianF8 handles lo,hi pairs

		for l < h {
			if ar[h] < pv {
				if ar[l] > pv {
					ar[l], ar[h] = ar[h], ar[l]
					h--
				}
				l++
			} else {
				if ar[l] <= pv { // extend ranges in balance
					l++
				}
				h--
			}
		}

		if l == h {
			if pv < ar[l] { // classify mid element
				h--
			} else {
				l++
			}
		}

		if h-lo < hi-l {
			h, hi = hi, h // [lo,h] is the longer range
			l, lo = lo, l
		}

		// branches below are optimally laid out for fewer jumps
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

		// range not long enough for new goroutine? max goroutines?
		// not atomic but good enough
		if hi-l < Mlr || ng >= Mxg {
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
