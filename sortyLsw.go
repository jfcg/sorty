/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "sync/atomic"

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
		if lsw(i, i-1, 0, 0) {
			return i
		}
	}
	return 0
}

// insertion sort ar[lo..hi], assumes lo < hi
func insertion(lsw Lesswap, lo, hi int) {

	for l, h := mid(lo, hi-1)-1, hi; l >= lo; {
		lsw(h, l, h, l)
		l--
		h--
	}
	for h := lo; ; {
		l := h
		h++
		k := h
		for lsw(k, l, k, l) {
			k = l
			l--
			if l < lo {
				break
			}
		}
		if h >= hi {
			break
		}
	}
}

// pivot divides ar[lo..hi] into 2n+1 equal intervals, sorts mid-points of them
// to find median-of-2n+1 pivot. ensures lo/hi ranges have at least n elements by
// moving 2n of mid-points to n positions at lo/hi ends.
// assumes n > 0, lo+4n+1 < hi. returns start,pivot,end for partitioning.
func pivot(lsw Lesswap, lo, hi, n int) (int, int, int) {
	m := mid(lo, hi)
	s := (hi - lo + 1) / (2*n + 1) // step > 1
	l, h := m-n*s, m+n*s

	for q, k := h, m-2*s; k >= l; { // insertion sort ar[m+i*s], i=-n..n
		lsw(q, k, q, k)
		q -= s
		k -= s
	}
	for r := l; ; {
		k := r
		r += s
		q := r
		for lsw(q, k, q, k) {
			q = k
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
	for {
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
		if l >= h {
			break
		}
	}
	if l == h && (h == pv || lsw(h, pv, 0, 0)) { // classify mid element
		l++
	}
	return l
}

// rearrange ar[l..a] and ar[b..h] into <= and >= pivot, assumes l <= a < pv < b <= h
// gap (a..b) expands until one of the intervals is fully consumed
func partition2(lsw Lesswap, l, a, pv, b, h int) (int, int) {
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

// new-goroutine partition
func gpart1(lsw Lesswap, l, pv, h int, ch chan int) {
	ch <- partition1(lsw, l, pv, h)
}

// concurrent dual partitioning
// returns m with ar[:m] <= pivot, ar[m:] >= pivot
func cdualpar(lsw Lesswap, lo, hi int, ch chan int) int {

	lo, pv, hi := pivot(lsw, lo, hi, 4) // median-of-9

	if hi-lo <= 2*Mlr { // guard against short remaining range
		return partition1(lsw, lo, pv, hi)
	}

	m := mid(lo, hi) // in pivot() lo/hi changed by possibly unequal amounts
	a, b := mid(lo, m), mid(m, hi)

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

// short range sort function, assumes Hmli <= hi-lo < Mlr
func short(lsw Lesswap, lo, hi int) {
start:
	l, pv, h := pivot(lsw, lo, hi, 2)
	l = partition1(lsw, l, pv, h) // median-of-5 partitioning
	h = l - 1
	no, n := h-lo, hi-l

	if no < n {
		n, no = no, n // [lo,hi] is the longer range
		l, lo = lo, l
	} else {
		h, hi = hi, h
	}

	if n >= Hmli {
		short(lsw, l, h) // recurse on the shorter range
		goto start
	}
	insertion(lsw, l, h) // at least one insertion range

	if no >= Hmli {
		goto start
	}
	insertion(lsw, lo, hi) // two insertion ranges
	return
}

// new-goroutine sort function
func glong(lsw Lesswap, lo, hi int, sv *syncVar) {
	long(lsw, lo, hi, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes hi-lo >= Mlr
func long(lsw Lesswap, lo, hi int, sv *syncVar) {
start:
	l, pv, h := pivot(lsw, lo, hi, 3)
	l = partition1(lsw, l, pv, h) // median-of-7 partitioning
	h = l - 1
	no, n := h-lo, hi-l

	if no < n {
		n, no = no, n // [lo,hi] is the longer range
		l, lo = lo, l
	} else {
		h, hi = hi, h
	}

	// branches below are optimal for fewer total jumps
	if n < Mlr { // at least one not-long range?

		if n >= Hmli {
			short(lsw, l, h)
		} else {
			insertion(lsw, l, h)
		}

		if no >= Mlr { // two not-long ranges?
			goto start
		}
		short(lsw, lo, hi) // we know no >= Hmli
		return
	}

	// single goroutine? max goroutines? not atomic but good enough
	if sv == nil || sv.ngr >= Mxg {
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

	n-- // high indice
	if n <= 2*Mlr || Mxg <= 1 {
		if n >= Mlr {
			long(lsw, 0, n, nil) // will not create goroutines or use ngr/done

		} else if n >= Hmli {
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
		// median-of-9 concurrent dual partitioning with done
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
		if n >= Mlr {
			if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
				panic("sorty: Sort: counter overflow")
			}
			go glong(lsw, l, h, &sv)

		} else if n >= Hmli {
			short(lsw, l, h)
		} else {
			insertion(lsw, l, h)
		}

		// longer range big enough? max goroutines?
		if no <= 2*Mlr || sv.ngr >= Mxg {
			break
		}
		// dual partition longer range
	}

	long(lsw, lo, hi, &sv) // we know hi-lo >= Mlr

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
