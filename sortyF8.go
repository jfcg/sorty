/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

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

// set pivot such that ar[l,l+1] <= pv <= ar[h-1,h]
func pivotF8(ar []float64, l, h int) (int, int, float64) {
	m := mid(l, h)
	vl, pv, vh, _ := slmhF8(ar[l], ar[m], ar[h])
	va, pv, vb, r := slmhF8(ar[l+1], pv, ar[h-1])

	if r > 0 && pv < vl {
		pv, vl = vl, pv
	}
	if r < 0 && vh < pv {
		vh, pv = pv, vh
	}
	ar[l], ar[l+1], ar[m], ar[h-1], ar[h] = vl, va, pv, vb, vh

	return l + 2, h - 2, pv
}

// partition ar into two groups: >= and <= pivot
func partitionF8(ar []float64, l, h int) (int, int) {
	l, h, pv := pivotF8(ar, l, h)
out:
	for ; l < h; l, h = l+1, h-1 {

		if ar[h] < pv { // avoid unnecessary comparisons
			for {
				if pv < ar[l] {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
				l++
				if l >= h {
					break out
				}
			}
		} else if pv < ar[l] { // extend ranges in balance
			for {
				h--
				if l >= h {
					break out
				}
				if ar[h] < pv {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
			}
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

// SortF8 concurrently sorts ar in ascending order.
func SortF8(ar []float64) {
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
		l, h := partitionF8(ar, lo, hi)

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

	arhi := len(ar) - 1
	if arhi > 2*Mlr {
		ng, done = 1, make(chan bool, 1)
		gsrt(0, arhi) // start master sort
		<-done
		return
	}

	if arhi >= Mli {
		srt(0, arhi) // single goroutine
		return
	}
	insertionF8(ar)
}
