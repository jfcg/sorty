package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// IsSortedU8 checks if ar is sorted in ascending order.
func IsSortedU8(ar []uint64) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

// insertion sort
func insertionU8(ar []uint64) {

	for l, h := len(ar)>>1-2, len(ar)-1; l >= 0; l, h = l-1, h-1 {
		if ar[h] < ar[l] {
			ar[l], ar[h] = ar[h], ar[l]
		}
	}

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
func slmhU8(vl, pv, vh uint64) (a, b, c uint64, r int) {
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
func medianU8(ar []uint64) uint64 {
	// lo, mid, hi
	h := len(ar) - 1
	m := h >> 1
	vl, va, pv, vb, vh := ar[0], ar[1], ar[m], ar[h-1], ar[h]

	vl, pv, vh, _ = slmhU8(vl, pv, vh)
	va, pv, vb, r := slmhU8(va, pv, vb)

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

// SortU8 concurrently sorts ar in ascending order.
func SortU8(ar []uint64) {
	var (
		ng        uint32         // number of sorting goroutines including this
		done      chan bool      // end signal
		srt, gsrt func(int, int) // recursive & new-goroutine sort functions
	)

	gsrt = func(lo, hi int) {
		srt(lo, hi)
		if atomic.AddUint32(&ng, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- false // we are the last, all done
		}
	}

	srt = func(lo, hi int) { // assumes hi-lo >= Mli
	start:
		l, h, pv := lo+2, hi-2, medianU8(ar[lo:hi+1]) // medianU8 handles lo,hi pairs

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
			insertionU8(ar[l : hi+1])

			if h-lo < Mli { // two short ranges?
				insertionU8(ar[lo : h+1])
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
			panic("SortU8: counter overflow")
		}
		go gsrt(lo, h) // start a new-goroutine sort on the longer range
		lo = l
		goto start
	}

	arhi := len(ar) - 1
	if arhi >= Mlr {
		ng, done = 1, make(chan bool, 1)
		gsrt(0, arhi) // start master sort
		<-done
		return
	}

	if arhi >= Mli {
		srt(0, arhi) // single goroutine
		return
	}
	insertionU8(ar)
}
