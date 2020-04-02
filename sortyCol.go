/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "sync/atomic"

// Lesser represents a collection of comparable elements.
type Lesser interface {
	// Len is the number of elements in the collection. First element has index 0.
	Len() int

	// Less (a strict ordering like < or >) reports if element-i should come
	// before element-k:
	//  Less(i,k) && Less(k,r) => Less(i,r)
	//  Less(i,k) => ! Less(k,i)
	Less(i, k int) bool
}

// Collection is same with standard library's sort.Interface. It represents a
// collection of sortable (as per Less) elements.
type Collection interface {
	Lesser

	// Swaps i-th and k-th elements.
	Swap(i, k int)
}

// IsSorted checks if ar is sorted.
func IsSorted(ar Lesser) bool {
	for i := ar.Len() - 1; i > 0; i-- {
		if ar.Less(i, i-1) {
			return false
		}
	}
	return true
}

// insertion sort
func insertion(ar Collection, lo, hi int) {

	for l, h := mid(lo, hi+1)-2, hi; l >= lo; l, h = l-1, h-1 {
		if ar.Less(h, l) {
			ar.Swap(h, l)
		}
	}

	for h := lo + 1; h <= hi; h++ {
		for l := h; ar.Less(l, l-1); {
			ar.Swap(l, l-1)
			l--
			if l <= lo {
				break
			}
		}
	}
}

// set such that ar[l,l+1] <= ar[m] = pivot <= ar[h-1,h]
func pivot(ar Collection, l, h int) (a, b, c int) {
	m := mid(l, h)
	if ar.Less(h, l) {
		ar.Swap(h, l)
	}
	if ar.Less(h, m) {
		ar.Swap(h, m)
	} else if ar.Less(m, l) {
		ar.Swap(m, l)
	}
	// ar[l] <= ar[m] <= ar[h]

	k, h := h, h-1
	if ar.Less(h, m) {
		k--
		ar.Swap(h, m)
		if ar.Less(m, l) {
			ar.Swap(m, l)
		}
	}
	l++

	if ar.Less(m, l) {
		ar.Swap(m, l)
		if k > h && ar.Less(h, h+1) {
			k = h
		}
		if ar.Less(k, m) {
			ar.Swap(k, m)
		}
	}

	return l + 1, h - 1, m
}

// partition ar into two groups: >= and <= pivot
func partition(ar Collection, l, h int) (int, int) {
	l, h, pv := pivot(ar, l, h)

	for {
		if ar.Less(h, pv) { // avoid unnecessary comparisons
			for {
				if ar.Less(pv, l) {
					ar.Swap(l, h)
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
				if ar.Less(h, pv) {
					ar.Swap(l, h)
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

// Sort concurrently sorts ar.
func Sort(ar Collection) {
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
		l, h := partition(ar, lo, hi)

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

		// range not long enough for new goroutine? max goroutines?
		// not atomic but good enough
		if hi-l < Mlr || ng >= Mxg {
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
	insertion(ar, 0, arhi)
}
