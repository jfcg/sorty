/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "sync/atomic"

// IsSorted checks if underlying collection of length n is sorted as per less().
// An existing lesswap() can be used like:
//  IsSorted(n, func(i, k int) bool { return lesswap(i, k, 0, 0) })
func IsSorted(n int, less func(i, k int) bool) bool {
	for i := n - 1; i > 0; i-- {
		if less(i, i-1) {
			return false
		}
	}
	return true
}

// Lesswap function operates on an underlying collection to be sorted as:
//  if less(i, k) { // strict ordering like < or >
//  	if r != s {
//  		swap(r, s)
//  	}
//  	return true
//  }
//  return false
type Lesswap func(i, k, r, s int) bool

// insertion sort
func insertion(lsw Lesswap, lo, hi int) {

	for l, h := mid(lo, hi+1)-2, hi; l >= lo; l, h = l-1, h-1 {
		lsw(h, l, h, l)
	}

	for h := lo + 1; h <= hi; h++ {
		for l := h; lsw(l, l-1, l, l-1); {
			l--
			if l <= lo {
				break
			}
		}
	}
}

// set such that ar[l,l+1] <= ar[m] = pivot <= ar[h-1,h]
func pivot(lsw Lesswap, l, h int) (int, int, int) {
	m := mid(l, h)
	lsw(h, l, h, l)
	if !lsw(h, m, h, m) {
		lsw(m, l, m, l)
	}
	// ar[l] <= ar[m] <= ar[h]

	k, h := h, h-1
	if lsw(h, m, h, m) {
		k--
		lsw(m, l, m, l)
	}

	l++
	if lsw(m, l, m, l) {
		if k > h && lsw(h, k, 0, 0) {
			k = h
		}
		lsw(k, m, k, m)
	}
	return l + 1, m, h - 1
}

// partition ar[l,h] into two groups: >= and <= pivot
func partition(lsw Lesswap, l, p, h int) int {
	for {
		if lsw(h, p, 0, 0) { // avoid unnecessary comparisons
			for {
				if lsw(p, l, h, l) {
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if lsw(p, l, 0, 0) { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l
				}
				if lsw(h, p, h, l) {
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

	if l == h && h != p && lsw(h, p, 0, 0) { // classify mid element
		l++
	}
	return l
}

// Sort concurrently sorts underlying collection of length n via lsw().
// Once for each non-trivial type you want to sort in a certain way, you
// can implement a custom sorting routine (for a slice for example) as:
//  func SortObjAsc(c []Obj) {
//  	lsw := func(i, k, r, s int) bool {
//  		if c[i].Key < c[k].Key { // your custom strict comparator (like < or >)
//  			if r != s {
//  				c[r], c[s] = c[s], c[r]
//  			}
//  			return true
//  		}
//  		return false
//  	}
//  	sorty.Sort(len(c), lsw)
//  }
func Sort(n int, lsw Lesswap) {
	var (
		mli       = Mli >> 1
		ngr       = uint32(1)    // number of sorting goroutines including this
		done      chan bool      // end signal
		srt, gsrt func(int, int) // recursive & new-goroutine sort functions
	)

	gsrt = func(lo, hi int) {
		srt(lo, hi)
		if atomic.AddUint32(&ngr, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- false // we are the last, all done
		}
	}

	srt = func(lo, hi int) { // assumes hi-lo >= mli
	start:
		l, p, h := pivot(lsw, lo, hi)
		l = partition(lsw, l, p, h)
		h = l - 1

		if h-lo < hi-l {
			h, hi = hi, h // [lo,h] is the longer range
			l, lo = lo, l
		}

		// branches below are optimally laid out for fewer jumps
		// at least one short range?
		if hi-l < mli {
			insertion(lsw, l, hi)

			if h-lo < mli { // two short ranges?
				insertion(lsw, lo, h)
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
			panic("Sort: counter overflow")
		}
		go gsrt(lo, h) // start a new-goroutine sort on the longer range
		lo = l
		goto start
	}

	n--
	if n > 2*Mlr {
		done = make(chan bool, 1)
		gsrt(0, n) // start master sort
		<-done
		return
	}

	if n >= mli {
		srt(0, n) // single goroutine
		return
	}
	insertion(lsw, 0, n)
}
