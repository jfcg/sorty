/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "sync/atomic"

// IsSortedF4 returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedF4(ar []float32) int {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return i
		}
	}
	return 0
}

// insertion sort, assumes 0 < hi < len(ar)
func insertionF4(ar []float32, hi int) {

	for l, h := (hi-3)>>1, hi; l >= 0; {
		if ar[h] < ar[l] {
			ar[l], ar[h] = ar[h], ar[l]
		}
		l--
		h--
	}
	for h := 0; ; {
		l := h
		h++
		v := ar[h]
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
		if h >= hi {
			break
		}
	}
}

// arrange median-of-5 as ar[l,l+1] <= ar[m] = pivot <= ar[h-1,h]
// This allows ar[l,l+1] and ar[h-1,h] to assist pivoting of the two sub-ranges in
// next pivot calls: 3 new values from a sub-range, 2 expectedly good values from
// parent range. Users of pivotF4() must ensure l+5 < h < len(ar)
func pivotF4(ar []float32, l, h int) (int, float32, int) {
	m := mid(l, h)
	vl, va, pv, vb, vh := ar[l], ar[l+1], ar[m], ar[h-1], ar[h]

	if vh < vl {
		vh, vl = vl, vh
	}
	if vb < va {
		vb, va = va, vb
	}
	if vh < vb {
		vh, vb = vb, vh
		vl, va = va, vl
	}

	if pv < vl {
		pv, vl = vl, pv
	}
	if vb < pv {
		vb, pv = pv, vb
		vl, va = va, vl
	}
	if pv < va {
		pv, va = va, pv
	}

	ar[l], ar[l+1], ar[m], ar[h-1], ar[h] = vl, va, pv, vb, vh
	return l + 2, pv, h - 2
}

// partition ar into >= and <= pivot, assumes l < h
func partitionF4(ar []float32, l, h int) int {
	l, pv, h := pivotF4(ar, l, h)
	for {
		if ar[h] < pv { // avoid unnecessary comparisons
			for {
				if pv < ar[l] {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if pv < ar[l] { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l
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

	if l == h && ar[h] < pv { // classify mid element
		l++
	}
	return l
}

// SortF4 concurrently sorts ar in ascending order.
func SortF4(ar []float32) {
	var (
		ngr  = uint32(1)    // number of sorting goroutines including this
		done chan bool      // end signal
		srt  func(int, int) // recursive sort function
	)

	gsrt := func(lo, no int) { // new-goroutine sort function
		srt(lo, no)
		if atomic.AddUint32(&ngr, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- false // we are the last, all done
		}
	}

	srt = func(lo, no int) { // assumes no >= Mli
	start:
		n := lo + no
		l := partitionF4(ar, lo, n)
		n -= l
		no -= n + 1

		if no < n {
			n, no = no, n // [lo,lo+no] is the longer range
			l, lo = lo, l
		}

		// branches below are optimally laid out for fewer jumps
		// at least one short range?
		if n < Mli {
			insertionF4(ar[l:], n)

			if no < Mli { // two short ranges?
				insertionF4(ar[lo:], no)
				return
			}
			goto start
		}

		// range not long enough for new goroutine? max goroutines?
		// not atomic but good enough
		if n < Mlr || ngr >= Mxg {
			srt(l, n) // start a recursive sort on the shorter range
			goto start
		}

		if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
			panic("SortF4: counter overflow")
		}
		go gsrt(lo, no) // start a new-goroutine sort on the longer range
		lo, no = l, n
		goto start
	}

	n := len(ar) - 1 // high indice
	if n > 2*Mlr {
		done = make(chan bool, 1)
		gsrt(0, n) // start master sort
		<-done
		return
	}
	if n >= Mli {
		srt(0, n) // single goroutine
		return
	}
	if n > 0 {
		insertionF4(ar, n) // length 2+
	}
}
