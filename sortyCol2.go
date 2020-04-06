/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "sync/atomic"

// Collection2 is an alternative for faster interface-based sorting
type Collection2 interface {
	Lesser

	// LessSwap must be equivalent to:
	//  if less(i, k) {
	//  	swap(r, s)
	//  	return true
	//  }
	//  return false
	LessSwap(i, k, r, s int) bool
}

func lsw(ar Collection2, a, b int) bool {
	return ar.LessSwap(a, b, a, b)
}

// insertion sort
func insertion2(ar Collection2, lo, hi int) {

	for l, h := mid(lo, hi+1)-2, hi; l >= lo; l, h = l-1, h-1 {
		lsw(ar, h, l)
	}

	for h := lo + 1; h <= hi; h++ {
		for l := h; lsw(ar, l, l-1); {
			l--
			if l <= lo {
				break
			}
		}
	}
}

// set such that ar[l,l+1] <= ar[m] = pivot <= ar[h-1,h]
func pivot2(ar Collection2, l, h int) (int, int, int) {
	m := mid(l, h)
	lsw(ar, h, l)
	if !lsw(ar, h, m) {
		lsw(ar, m, l)
	}
	// ar[l] <= ar[m] <= ar[h]

	k, h := h, h-1
	if lsw(ar, h, m) {
		k--
		lsw(ar, m, l)
	}

	l++
	if lsw(ar, m, l) {
		if k > h && ar.Less(h, k) {
			k = h
		}
		lsw(ar, k, m)
	}
	return l + 1, h - 1, m
}

// partition ar into two groups: >= and <= pivot
func partition2(ar Collection2, l, h int) (int, int) {
	l, h, pv := pivot2(ar, l, h)

	for {
		if ar.Less(h, pv) { // avoid unnecessary comparisons
			for {
				if ar.LessSwap(pv, l, h, l) {
					break
				}
				l++
				if l >= h {
					return l + 1, h
				}
			}
		} else if ar.Less(pv, l) { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l, h - 1
				}
				if ar.LessSwap(h, pv, h, l) {
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
		if ar.Less(pv, l) { // classify mid element
			h--
		} else {
			l++
		}
	}
	return l, h
}

// Sort2 concurrently sorts ar.
func Sort2(ar Collection2) {
	var (
		arhi, mli = ar.Len() - 1, Mli >> 2
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

	srt = func(lo, hi int) { // assumes hi-lo >= mli
	start:
		l, h := partition2(ar, lo, hi)

		if h-lo < hi-l {
			h, hi = hi, h // [lo,h] is the longer range
			l, lo = lo, l
		}

		// branches below are optimally laid out for fewer jumps
		// at least one short range?
		if hi-l < mli {
			insertion2(ar, l, hi)

			if h-lo < mli { // two short ranges?
				insertion2(ar, lo, h)
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
			panic("Sort2: counter overflow")
		}
		go gsrt(lo, h) // start a new-goroutine sort on the longer range
		lo = l
		goto start
	}

	if arhi > 2*Mlr {
		ng, done = 1, make(chan bool, 1)
		gsrt(0, arhi) // start master sort
		<-done
		return
	}

	if arhi >= mli {
		srt(0, arhi) // single goroutine
		return
	}
	insertion2(ar, 0, arhi)
}
