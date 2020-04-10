/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "sync/atomic"

// IsSortedU4 checks if ar is sorted in ascending order.
func IsSortedU4(ar []uint32) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

// insertion sort
func insertionU4(ar []uint32) {

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

// set such that ar[l,l+1] <= ar[m] = pivot <= ar[h-1,h]
func pivotU4(ar []uint32, l, h int) (int, uint32, int) {
	m := mid(l, h)
	vl, va, pv, vb, vh := ar[l], ar[l+1], ar[m], ar[h-1], ar[h]

	if vh < vl {
		vh, vl = vl, vh
	}
	if vh < pv {
		vh, pv = pv, vh
	} else if pv < vl {
		pv, vl = vl, pv
	}
	// vl <= pv <= vh

	if vb < pv {
		vb, pv = pv, vb
		if pv < vl {
			pv, vl = vl, pv
		}
	}

	if pv < va {
		pv, va = va, pv
		if vh < vb {
			vh, vb = vb, vh
		}
		if vb < pv {
			vb, pv = pv, vb
		}
	}

	ar[l], ar[l+1], ar[m], ar[h-1], ar[h] = vl, va, pv, vb, vh
	return l + 2, pv, h - 2
}

// partition ar into two groups: >= and <= pivot
func partitionU4(ar []uint32, l, h int) (int, int) {
	l, pv, h := pivotU4(ar, l, h)

	for {
		if ar[h] < pv { // avoid unnecessary comparisons
			for {
				if pv < ar[l] {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
				l++
				if l >= h {
					return l + 1, h
				}
			}
		} else if pv < ar[l] { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l, h - 1
				}
				if ar[h] < pv {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
			}
		}
		l++
		h--
		if l >= h {
			break
		}
	}

	if l == h {
		if pv < ar[l] { // classify mid element
			h--
		} else {
			l++
		}
	}
	return l, h
}

// SortU4 concurrently sorts ar in ascending order.
func SortU4(ar []uint32) {
	var (
		ngr       uint32         // number of sorting goroutines including this
		done      chan bool      // end signal
		srt, gsrt func(int, int) // recursive & new-goroutine sort functions
	)

	gsrt = func(lo, hi int) {
		srt(lo, hi)
		if atomic.AddUint32(&ngr, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- false // we are the last, all done
		}
	}

	srt = func(lo, hi int) { // assumes hi-lo >= Mli
	start:
		l, h := partitionU4(ar, lo, hi)

		if h-lo < hi-l {
			h, hi = hi, h // [lo,h] is the longer range
			l, lo = lo, l
		}

		// branches below are optimally laid out for fewer jumps
		// at least one short range?
		if hi-l < Mli {
			insertionU4(ar[l : hi+1])

			if h-lo < Mli { // two short ranges?
				insertionU4(ar[lo : h+1])
				return
			}
			hi = h
			goto start
		}

		// range not long enough for new goroutine? max goroutines?
		// not atomic but good enough
		if hi-l < Mlr || ngr >= Mxg {
			srt(l, hi) // start a recursive sort on the shorter range
			hi = h
			goto start
		}

		if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
			panic("SortU4: counter overflow")
		}
		go gsrt(lo, h) // start a new-goroutine sort on the longer range
		lo = l
		goto start
	}

	arhi := len(ar) - 1
	if arhi > 2*Mlr {
		ngr, done = 1, make(chan bool, 1)
		gsrt(0, arhi) // start master sort
		<-done
		return
	}

	if arhi >= Mli {
		srt(0, arhi) // single goroutine
		return
	}
	insertionU4(ar)
}
