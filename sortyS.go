package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// IsSortedS checks if ar is sorted in ascending order.
func IsSortedS(ar []string) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

// insertion sort
func insertionS(ar []string) {
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
func slmhS(vl, pv, vh string) (a, b, c string, r int) {
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
func medianS(ar []string) string {
	// lo, mid, hi
	h := len(ar) - 1
	m := h >> 1
	vl, va, pv, vb, vh := ar[0], ar[1], ar[m], ar[h-1], ar[h]

	vl, pv, vh, _ = slmhS(vl, pv, vh)
	va, pv, vb, r := slmhS(va, pv, vb)

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

// SortS concurrently sorts ar in ascending order.
func SortS(ar []string) {
	if len(ar) <= Mli {
		insertionS(ar)
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
		l, h = lo+2, hi-2 // medianS handles lo,hi pairs

		for pv := medianS(ar[lo : hi+1]); l <= h; {
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

		if h-lo < hi-l {
			h, hi = hi, h // [lo,h] is the longer range
			l, lo = lo, l
		}

		// branches below are optimally laid out for less # of jumps
		// at least one short range?
		if hi-l < Mli {
			insertionS(ar[l : hi+1])

			if h-lo < Mli { // two short ranges?
				insertionS(ar[lo : h+1])
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
			panic("SortS: counter overflow")
		}
		go gsrt(lo, h) // start a new-goroutine sort on the longer range
		lo = l
		goto start
	}

	gsrt(0, len(ar)-1) // start master sort
	<-done
}
