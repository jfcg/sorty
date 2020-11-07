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

// insertion sort ar[lo..lo+no], assumes no > 0
func insertion(lsw Lesswap, lo, no int) {
	hi := lo + no
	for l, h := mid(lo, hi-1)-1, hi; l >= lo; {
		lsw(h, l, h, l)
		l--
		h--
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

// pivot divides [lo..hi] range into 2n+1 equal intervals, sorts mid-points of them
// to find median-of-2n+1 pivot. ensures lo/hi ranges have at least n elements by
// moving 2n of mid-points to n positions at lo/hi ends.
// assumes n > 0, lo+4n+1 < hi. returns start,pivot,end for partitioning.
func pivot(lsw Lesswap, lo, hi, n int) (int, int, int) {
	m := mid(lo, hi)
	s := int(uint(hi-lo+1) / uint(2*n+1)) // step > 1
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

// concurrent dual partitioning
// returns m with ar[:m] <= pivot, ar[m:] >= pivot
func cdualpar(par chan int, lsw Lesswap, lo, hi int) int {

	lo, pv, hi := pivot(lsw, lo, hi, 4) // median-of-9

	if hi-lo <= 2*Mlr { // guard against short remaining range
		return partition1(lsw, lo, pv, hi)
	}

	m := mid(lo, hi) // in pivot() lo/hi changed by possibly unequal amounts
	a, b := mid(lo, m), mid(m, hi)

	go func(l, h int) {
		par <- partition1(lsw, l, pv, h) // mid half range
	}(a, b)

	a, b = partition2(lsw, lo, a-1, pv, b+1, hi) // left/right quarter ranges
	m = <-par

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

// short range sort function, assumes Hmli <= no < Mlr
func short(lsw Lesswap, lo, no int) {
start:
	n := lo + no
	l, pv, h := pivot(lsw, lo, n, 2) // median-of-5
	l = partition1(lsw, l, pv, h)
	n -= l
	no -= n + 1

	if no < n {
		n, no = no, n // [lo,lo+no] is the longer range
		l, lo = lo, l
	}

	if n >= Hmli {
		short(lsw, l, n) // recurse on the shorter range
		goto start
	}
	insertion(lsw, l, n) // at least one insertion range

	if no >= Hmli {
		goto start
	}
	insertion(lsw, lo, no) // two insertion ranges
	return
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
		ngr  = uint32(1)    // number of sorting goroutines including this
		done chan int       // end signal
		long func(int, int) // long range sort function
	)

	glong := func(lo, no int) { // new-goroutine sort function
		long(lo, no)
		if atomic.AddUint32(&ngr, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- 0 // we are the last, all done
		}
	}

	long = func(lo, no int) { // assumes no >= Mlr
	start:
		n := lo + no
		l, pv, h := pivot(lsw, lo, n, 3) // median-of-7
		l = partition1(lsw, l, pv, h)
		n -= l
		no -= n + 1

		if no < n {
			n, no = no, n // [lo,lo+no] is the longer range
			l, lo = lo, l
		}

		// branches below are optimal for fewer total jumps
		if n < Mlr { // at least one not-long range?
			if n >= Hmli {
				short(lsw, l, n)
			} else {
				insertion(lsw, l, n)
			}

			if no >= Mlr { // two not-long ranges?
				goto start
			}
			short(lsw, lo, no) // we know no >= Hmli
			return
		}

		// max goroutines? not atomic but good enough
		if ngr >= Mxg {
			long(l, n) // recurse on the shorter range
			goto start
		}

		if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
			panic("Sort: long: counter overflow")
		}
		// new-goroutine sort on the longer range only when
		// both ranges are big and max goroutines is not exceeded
		go glong(lo, no)
		lo, no = l, n
		goto start
	}

	n-- // high indice
	if n <= 2*Mlr {
		if n >= Mlr {
			long(0, n) // will not create goroutines or use ngr/done
		} else if n >= Hmli {
			short(lsw, 0, n)
		} else if n > 0 {
			insertion(lsw, 0, n)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	done = make(chan int, 1) // maybe this goroutine will be the last
	lo, no := 0, n
	for {
		// concurrent dual partitioning with done
		l := cdualpar(done, lsw, lo, n)
		n -= l
		no -= n + 1

		if no < n {
			n, no = no, n // [lo,lo+no] is the longer range
			l, lo = lo, l
		}

		// handle shorter range
		if n >= Mlr {
			if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
				panic("Sort: dual: counter overflow")
			}
			go glong(l, n)

		} else if n >= Hmli {
			short(lsw, l, n)
		} else {
			insertion(lsw, l, n)
		}

		// longer range big enough? max goroutines?
		if no <= 2*Mlr || ngr >= Mxg {
			break
		}
		n = lo + no // dual partition longer range
	}

	glong(lo, no) // we know no >= Mlr
	<-done
}
