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

// insertion sort ar[lo..hi], assumes lo+2 < hi
func insertion(lsw Lesswap, lo, hi int) {

	for l, h := mid(lo, hi-1)-1, hi; ; {
		lsw(h, l, h, l)
		h--
		l--
		if l < lo {
			break
		}
	}
	for h := lo; ; {
		for l := h; lsw(l+1, l, l+1, l); {
			l--
			if l < lo {
				break
			}
		}
		h++
		if h >= hi {
			break
		}
	}
}

// arrange median-of-9 as ar[l,l+1, a-1,a] <= ar[m] = pivot <= ar[b,b+1, h-1,h]
// if dual: a,b = mid(l,m), mid(m,h) else: a,b = l+3,h-3
// pivot() ensures partitioning yields ranges of length 4+
// users of pivot() must ensure l+11 < h
func pivot(lsw Lesswap, l, h int, dual bool) (int, int, int) {
	s := [9]int{l, l + 1, 0, 0, mid(l, h), 0, 0, h - 1, h}
	if dual {
		s[3], s[5] = mid(l, s[4]), mid(s[4], h)
	} else {
		s[3], s[5] = l+3, h-3
	}
	s[2], s[6] = s[3]-1, s[5]+1

	for i := 2; i >= 0; i-- { // insertion sort via s
		lsw(s[i+6], s[i], s[i+6], s[i])
	}
	for i := 1; i < len(s); i++ {
		for k, r := i-1, s[i]; lsw(r, s[k], r, s[k]); {
			r = s[k]
			k--
			if k < 0 {
				break
			}
		}
	}
	return s[3] + 1, s[4], s[5] - 1 // l,pv,h suitable for partition()
}

// partition ar[l..h] into <= and >= pivot, assumes l < pv < h
func partition(lsw Lesswap, l, pv, h int) int {
	for {
		if lsw(h, pv, 0, 0) { // avoid unnecessary comparisons
			for {
				if lsw(pv, l, h, l) {
					break
				}
				l++
				if l >= pv { // until pv & avoid it
					l++
					goto next
				}
			}
		} else if lsw(pv, l, 0, 0) { // extend ranges in balance
			for {
				h--
				if pv >= h { // until pv & avoid it
					h--
					goto next
				}
				if lsw(h, pv, h, l) {
					break
				}
			}
		}
		l++
		h--
		if l >= pv {
			l++
			break
		}
		if pv >= h {
			h--
			goto next
		}
	}
	if pv >= h {
		h--
	}

next:
	for l < h {
		if lsw(h, pv, 0, 0) { // avoid unnecessary comparisons
			for {
				if lsw(pv, l, h, l) {
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if lsw(pv, l, 0, 0) { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l
				}
				if lsw(h, pv, h, l) {
					break
				}
			}
		}
		l++
		h--
	}
	if l == h && lsw(h, pv, 0, 0) { // classify mid element
		l++
	}
	return l
}

// rearrange ar[l..a] & ar[b..h] into <= and >= pivot, assumes l <= a < pv < b <= h
// gap (a..b) expands until one of the intervals is fully consumed
func dpartition(lsw Lesswap, l, a, pv, b, h int) (int, int) {
	for {
		if lsw(b, pv, 0, 0) { // avoid unnecessary comparisons
			for {
				if lsw(pv, a, b, a) {
					break
				}
				a--
				if a < l {
					return a, b
				}
			}
		} else if lsw(pv, a, 0, 0) { // extend ranges in balance
			for {
				b++
				if b > h {
					return a, b
				}
				if lsw(b, pv, b, a) {
					break
				}
			}
		}
		a--
		b++
		if a < l || b > h {
			return a, b
		}
	}
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
		mli  = Mli >> 1
		ngr  = uint32(1)               // number of sorting goroutines including this
		done chan bool                 // end signal
		srt  func(*chan int, int, int) // recursive sort function
	)

	// new-goroutine sort function
	gsrt := func(lo, hi int) {
		var par chan int // dual partitioning channel
		srt(&par, lo, hi)

		if atomic.AddUint32(&ngr, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- false // we are the last, all done
		}
	}

	// dual partitioning
	dualpar := func(par *chan int, lo, hi int) int {

		if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
			panic("Sort: dualpar: counter overflow")
		}
		if *par == nil {
			*par = make(chan int) // make it when 1st time we need it
		}
		a, pv, b := pivot(lsw, lo, hi, true)

		go func(a, b int) {
			*par <- partition(lsw, a, pv, b)

			if atomic.AddUint32(&ngr, ^uint32(0)) == 0 { // decrease goroutine counter
				panic("Sort: dualpar: counter underflow")
			}
		}(a, b)

		a -= 3
		b += 3
		lo += 2
		hi -= 2
		a, b = dpartition(lsw, lo, a, pv, b, hi)
		m := <-*par // <= and >= boundary

		// only one gap is possible
		for ; lo <= a; a-- { // gap left in low range?
			if lsw(pv, a, m-1, a) {
				m--
				if m == pv { // swapped pivot when closing gap?
					pv = a // Thanks to my wife Tansu who discovered this
				}
			}
		}
		for ; b <= hi; b++ { // gap left in high range?
			if lsw(b, pv, b, m) {
				if m == pv { // swapped pivot when closing gap?
					pv = b // It took days of agony to discover these two if's :D
				}
				m++
			}
		}
		return m
	}

	srt = func(par *chan int, lo, hi int) { // assumes hi-lo >= mli
	start:
		var l int
		// range long enough for dual partitioning? available goroutine?
		// not atomic but good enough
		if hi-lo > 2*Mlr && ngr < Mxg {
			l = dualpar(par, lo, hi)
		} else {
			a, pv, b := pivot(lsw, lo, hi, false)
			l = partition(lsw, a, pv, b)
		}
		h := l - 1

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
			srt(par, l, hi) // start a recursive sort on the shorter range
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

	n-- // high indice
	if n > 2*Mlr {
		done = make(chan bool, 1)
		gsrt(0, n) // start master sort
		<-done
		return
	}
	if n >= mli {
		srt(nil, 0, n) // single goroutine
		return
	}
	if n > 2 {
		insertion(lsw, 0, n) // length 4+
		return
	}

	if n > 0 { // handle arrays of length 2,3
		for {
			lsw(1, 0, 1, 0)
			n--
			if n <= 0 || !lsw(2, 1, 2, 1) {
				return
			}
		}
	}
}
