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

// arrange median-of-5 as ar[l,l+1] <= ar[m] = pivot <= ar[h-1,h]
// This allows ar[l,l+1] and ar[h-1,h] to assist pivoting of the two sub-ranges in
// next pivot5() calls: 3 new values from a sub-range, 2 expectedly good values from
// parent range. Users of pivot5() must ensure l+5 < h
func pivot5(lsw Lesswap, l, h int) (int, int, int) {
	e := l
	lsw(h, l, h, l)
	l++
	h--
	b, d := h, l
	lsw(h, l, h, l)

	if lsw(h+1, h, 0, 0) {
		d, e = e, d
		b++
	}
	c := mid(l, h)
	lsw(c, e, c, e)

	if lsw(b, c, b, c) {
		d = e
	}
	lsw(c, d, c, d)

	return l + 1, c, h - 1 // l,pv,h suitable for part1()
}

// arrange median-of-9 as ar[a-1] <= ar[l,l+1] <= ar[a] <= ar[m] = pivot <= ar[b] <=
// ar[h-1,h] <= ar[b+1] where m,a,b = mid(l,h), mid(l,m), mid(m,h). After pivot9() ar
// will look like: 2nd 3rd .. 1st 4th .. 5th=pivot .. 6th 9th .. 7th 8th.
// This allows ar[l,l+1] and ar[h-1,h] to assist pivoting of the two sub-ranges in
// next pivot9() or pivot5() calls: 7 or 3 new values from a sub-range, 2 expectedly
// good values from parent range. Users of pivot9() must ensure l+11 < h
func pivot9(lsw Lesswap, l, h int) (int, int, int) {

	s := [9]int{0, l, l + 1, 0, mid(l, h), 0, h - 1, h, 0}
	s[3], s[5] = mid(l, s[4]), mid(s[4], h)
	s[0], s[8] = s[3]-1, s[5]+1

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
	return s[0] - 1, s[4], s[8] + 1 // a,pv,b suitable for part2(), part1s()
}

// partition ar[l..h] into <= and >= pivot, assumes l < pv < h
func part1(lsw Lesswap, l, pv, h int) int {
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
	return part0(lsw, l, pv, h)
}

// partition ar[l..h] into <= and >= pivot, assumes pv is outside [l..h]
func part0(lsw Lesswap, l, pv, h int) int {
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

// rearrange ar[l..a] and ar[b..h] into <= and >= pivot, assumes l <= a < pv < b <= h
// gap (a..b) expands until one of the intervals is fully consumed
func part2(lsw Lesswap, l, a, pv, b, h int) (int, int) {
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
func cdualpar(par chan int, lsw Lesswap, lo, hi int) int {
	a, pv, b := pivot9(lsw, lo, hi)

	go func(l, h int) {
		par <- part1(lsw, l, pv, h)
	}(a+3, b-3)

	lo += 2
	hi -= 2
	a, b = part2(lsw, lo, a, pv, b, hi)
	m := <-par

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

// partition ar[l..h] into <= and >= pivot with pivot9(), skipping around a,b
// assumes l+11 < h
func part1s(lsw Lesswap, l, h int) int {
	a, pv, b := pivot9(lsw, l, h)
	a++
	b-- // skipping a,a+1,b-1,b
	l += 2
	h -= 2
	for {
		if lsw(h, pv, 0, 0) { // avoid unnecessary comparisons
			for {
				if lsw(pv, l, h, l) {
					break
				}
				l++
				if l >= a { // until a & avoid pair
					l++
					if a == pv {
						goto next
					}
					l++
					a = pv
				}
			}
		} else if lsw(pv, l, 0, 0) { // extend ranges in balance
			for {
				h--
				if b >= h { // until b & avoid pair
					h--
					if b == pv {
						goto next
					}
					h--
					b = pv
				}
				if lsw(h, pv, h, l) {
					break
				}
			}
		}
		l++
		h--
		if l >= a {
			l++
			if a == pv {
				break
			}
			l++
			a = pv
		}
		if b >= h {
			h--
			if b == pv {
				goto next
			}
			h--
			b = pv
		}
	}
	if b >= h {
		h--
		if b != pv {
			h--
		}
	}

next:
	return part0(lsw, l, pv, h)
}

// short range sort function, assumes Hmli <= no <= 2*Mlr
func short(lsw Lesswap, lo, no int) {
start:
	n := lo + no
	l, pv, h := pivot5(lsw, lo, n)
	l = part1(lsw, l, pv, h) // partitioning
	n -= l
	no -= n + 1

	if no < n {
		n, no = no, n // [lo,lo+no] is the longer range
		l, lo = lo, l
	}

	// branches below are optimal for fewer total jumps
	if n < Hmli { // at least one insertion range?
		insertion(lsw, l, l+n)

		if no >= Hmli { // two insertion ranges?
			goto start
		}
		insertion(lsw, lo, lo+no)
		return
	}

	short(lsw, l, n) // start a recursive sort on the shorter range
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
		l := part1s(lsw, lo, n) // partitioning with pivot9()
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
				insertion(lsw, l, l+n)
			}

			if no >= Mlr { // two not-long ranges?
				goto start
			}
			short(lsw, lo, no) // we know no >= Hmli
			return
		}

		// max goroutines? not atomic but good enough
		if ngr >= Mxg {
			long(l, n) // start a recursive sort on the shorter range
			goto start
		}

		if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
			panic("Sort: long: counter overflow")
		}
		go glong(lo, no) // start a new-goroutine sort on the longer range
		lo, no = l, n
		goto start
	}

	n-- // high indice
	if n <= 2*Mlr {
		if n >= Hmli {
			short(lsw, 0, n) // single goroutine
		} else if n > 0 {
			insertion(lsw, 0, n)
		}
		return
	}

	done = make(chan int, 1)
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
			insertion(lsw, l, l+n)
		}

		if no <= 2*Mlr || ngr >= Mxg {
			break
		}
		n = lo + no // dual partition longer range
	}

	glong(lo, no) // we know no >= Mlr
	<-done
}
