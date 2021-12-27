/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"sync/atomic"

	"github.com/jfcg/sixb"
)

// Lesswap function operates on an underlying collection to be sorted as:
//  if less(i, k) { // strict ordering like < or >
//  	if r != s {
//  		swap(r, s)
//  	}
//  	return true
//  }
//  return false
type Lesswap func(i, k, r, s int) bool

// IsSorted returns 0 if underlying collection of length n is sorted,
// otherwise it returns i > 0 with less(i,i-1) = true.
func IsSorted(n int, lsw Lesswap) int {
	for i := n - 1; i > 0; i-- {
		if lsw(i, i-1, i, i) { // 3rd=4th disables swap
			return i
		}
	}
	return 0
}

// insertion sort ar[l..h], assumes l < h
func insertion(lsw Lesswap, l, h int) {

	for i, k := l, l+(MaxLenInsFC+1)/3; k <= h; { // pre-sort
		lsw(k, i, k, i)
		i++
		k++
	}
	for k := l; ; { // insertion
		i := k
		k++
		q := k
		for lsw(q, i, q, i) {
			q--
			i--
			if i < l {
				break
			}
		}
		if k >= h {
			break
		}
	}
}

// pivot selects 2n+1 equidistant samples from ar[lo..hi] that minimizes max distance to
// any non-selected member, calculates median-of-2n+1 pivot from samples. ensures lo/hi
// ranges have at least n elements by moving sorted samples to n positions at lo/hi ends.
// assumes n > 0, lo+4n+1 < hi. returns start,pivot,end for partitioning.
func pivot(lsw Lesswap, lo, hi, n int) (int, int, int) {

	m := sixb.MeanI(lo, hi)
	s := (hi - lo + 1) / (2*n + 1) // step > 1
	l, h := m-n*s, m+n*s
	if l-lo >= n && hi-h > (s+1)>>1 {
		s++
		l -= n
		h += n
	}

	lsw(h, l, h, l)
	for r := l; ; { // insertion sort ar[m+i*s], i=-n..n
		k := r
		r += s
		q := r
		for lsw(q, k, q, k) {
			q -= s
			k -= s
			if k < l {
				break
			}
		}
		if r >= h {
			break
		}
	}

	// move hi mid-points to hi end
	for {
		if h == hi || lsw(hi, h, hi, h) {
			h -= s
		}
		hi--
		if h <= m {
			break
		}
	}

	// move lo mid-points to lo end
	for {
		if l == lo || lsw(l, lo, l, lo) {
			l += s
		}
		lo++
		if l >= m {
			break
		}
	}
	return lo, m, hi // lo <= m-s+1, m+s-1 <= hi
}

// partition ar[l..h] into <= and >= pivot, assumes l < h
// returns m with ar[:m] <= pivot, ar[m:] >= pivot
func partition1(lsw Lesswap, l, pv, h int) int {
	// avoid unnecessary comparisons, extend ranges in balance
	for {
		if lsw(h, pv, h, h) { // 3rd=4th disables swap
			for {
				if lsw(pv, l, h, l) {
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if lsw(pv, l, l, l) { // 3rd=4th disables swap
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
		if l >= h {
			break
		}
	}
	// classify mid element
	if l == h && h != pv && lsw(h, pv, h, h) { // 3rd=4th disables swap
		l++
	}
	return l
}

// rearrange ar[l..a] and ar[b..h] into <= and >= pivot, assumes l <= a < pv < b <= h
// gap (a..b) expands until one of the intervals is fully consumed
func partition2(lsw Lesswap, l, a, pv, b, h int) (int, int) {
	// avoid unnecessary comparisons, extend ranges in balance
	for {
		if lsw(b, pv, b, b) { // 3rd=4th disables swap
			for {
				if lsw(pv, a, b, a) {
					break
				}
				a--
				if a < l {
					return a, b
				}
			}
		} else if lsw(pv, a, a, a) { // 3rd=4th disables swap
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

// new-goroutine partition
func gpart1(lsw Lesswap, l, pv, h int, ch chan int) {
	ch <- partition1(lsw, l, pv, h)
}

// concurrent dual partitioning
// returns m with ar[:m] <= pivot, ar[m:] >= pivot
func cdualpar(lsw Lesswap, lo, hi int, ch chan int) int {

	lo, pv, hi := pivot(lsw, lo, hi, 4) // median-of-9 pivot

	if hi-lo <= 2*MaxLenRec { // guard against short remaining range
		return partition1(lsw, lo, pv, hi)
	}

	m := sixb.MeanI(lo, hi) // in pivot() lo/hi changed by possibly unequal amounts
	a, b := sixb.MeanI(lo, m), sixb.MeanI(m, hi)

	go gpart1(lsw, a+1, pv, b-1, ch) // mid half range

	a, b = partition2(lsw, lo, a, pv, b, hi) // left/right quarter ranges
	m = <-ch

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

// short range sort function, assumes MaxLenInsFC <= hi-lo < MaxLenRec
func short(lsw Lesswap, lo, hi int) {
start:
	l, pv, h := pivot(lsw, lo, hi, 2) // median-of-5 pivot
	l = partition1(lsw, l, pv, h)
	h = l - 1
	no, n := h-lo, hi-l

	if no < n {
		n, no = no, n // [lo,hi] is the longer range
		l, lo = lo, l
	} else {
		h, hi = hi, h
	}

	if n >= MaxLenInsFC {
		short(lsw, l, h) // recurse on the shorter range
		goto start
	}
	// at least one insertion range, presort+insertion inlined
psort:
	for i, k := l, l+(MaxLenInsFC+1)/3; k <= h; { // pre-sort
		lsw(k, i, k, i)
		i++
		k++
	}
	for k := l; ; { // insertion
		i := k
		k++
		q := k
		for lsw(q, i, q, i) {
			q--
			i--
			if i < l {
				break
			}
		}
		if k >= h {
			break
		}
	}

	if no >= MaxLenInsFC {
		goto start
	}
	if lo != l {
		l, h = lo, hi
		goto psort // two insertion ranges
	}
}

// new-goroutine sort function
func glong(lsw Lesswap, lo, hi int, sv *syncVar) {
	long(lsw, lo, hi, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes hi-lo >= MaxLenRec
func long(lsw Lesswap, lo, hi int, sv *syncVar) {
start:
	l, pv, h := pivot(lsw, lo, hi, 3) // median-of-7 pivot
	l = partition1(lsw, l, pv, h)
	h = l - 1
	no, n := h-lo, hi-l

	if no < n {
		n, no = no, n // [lo,hi] is the longer range
		l, lo = lo, l
	} else {
		h, hi = hi, h
	}

	// branches below are optimal for fewer total jumps
	if n < MaxLenRec { // at least one not-long range?

		if n >= MaxLenInsFC {
			short(lsw, l, h)
		} else {
			insertion(lsw, l, h)
		}

		if no >= MaxLenRec { // two not-long ranges?
			goto start
		}
		short(lsw, lo, hi) // we know no >= MaxLenInsFC
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		long(lsw, l, h, sv) // recurse on the shorter range
		goto start
	}

	if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
		panic("sorty: long: counter overflow")
	}
	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	go glong(lsw, lo, hi, sv)
	lo, hi = l, h
	goto start
}

// Sort concurrently sorts underlying collection of length n via lsw().
// Once for each non-trivial type you want to sort in a certain way, you
// can implement a custom sorting routine (for a slice for example) as:
//  func SortObjAsc(c []Obj) {
//  	lsw := func(i, k, r, s int) bool {
//  		if c[i].Key < c[k].Key { // strict comparator like < or >
//  			if r != s {
//  				c[r], c[s] = c[s], c[r]
//  			}
//  			return true
//  		}
//  		return false
//  	}
//  	sorty.Sort(len(c), lsw)
//  }
// Lesswap is a 'contract' between users and sorty library. Strict
// comparator, r!=s check, swap and returns are all strictly necessary.
func Sort(n int, lsw Lesswap) {

	n-- // high indice
	if n <= 2*MaxLenRec || MaxGor <= 1 {

		if n >= MaxLenRec { // single-goroutine sorting
			long(lsw, 0, n, nil)
		} else if n >= MaxLenInsFC {
			short(lsw, 0, n)
		} else if n > 0 {
			insertion(lsw, 0, n)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	lo, hi := 0, n
	for {
		// concurrent dual partitioning with done
		l := cdualpar(lsw, lo, hi, sv.done)
		h := l - 1
		no, n := h-lo, hi-l

		if no < n {
			n, no = no, n // [lo,hi] is the longer range
			l, lo = lo, l
		} else {
			h, hi = hi, h
		}

		// handle shorter range
		if n >= MaxLenRec {
			if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
				panic("sorty: Sort: counter overflow")
			}
			go glong(lsw, l, h, &sv)

		} else if n >= MaxLenInsFC {
			short(lsw, l, h)
		} else {
			insertion(lsw, l, h)
		}

		// longer range big enough? max goroutines?
		if no <= 2*MaxLenRec || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	long(lsw, lo, hi, &sv) // we know hi-lo >= MaxLenRec

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
