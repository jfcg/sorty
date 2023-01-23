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
//
//	if less(i, k) { // strict ordering like < or >
//		if r != s {
//			swap(r, s)
//		}
//		return true
//	}
//	return false
type Lesswap func(i, k, r, s int) bool

// IsSorted returns 0 if underlying collection of length n is sorted,
// otherwise it returns i > 0 with less(i,i-1) = true.
//
//go:nosplit
func IsSorted(n int, lsw Lesswap) int {
	for i := n - 1; i > 0; i-- {
		if lsw(i, i-1, i, i) { // 3rd=4th disables swap
			return i
		}
	}
	return 0
}

// insertion sort ar[lo..hi]
//
//go:nosplit
func insertion(lsw Lesswap, lo, hi int) {
	for h := lo + 1; h <= hi; h++ {
		for l := h; lsw(l, l-1, l, l-1); {
			l--
			if l <= lo {
				break
			}
		}
	}
}

// pivot selects n equidistant samples from slc[lo:hi+1] that minimizes max distance
// to non-selected members, then calculates median-of-n pivot from samples.
// Assumes odd n ≥ 3 and len(slc) ≥ 2n. Returns pivot position.
// Moves one sorted sample to each end to ensure sub-slices have lengths ≥ 1
//
//go:nosplit
func pivot(lsw Lesswap, lo, hi int, n uint) int {

	f, s, l := minMaxSample(uint(hi+1-lo), n)
	first := lo + int(f)
	step := int(s)
	last := lo + int(l)

	// insertion sort slc[first + j * step], j=0,1,..
	for h := first + step; h <= last; h += step {
		for l := h; lsw(l, l-step, l, l-step); {
			l -= step
			if l <= first {
				break
			}
		}
	}

	// move one sorted sample to each end
	lsw(first, lo, first, lo)
	lsw(hi, last, hi, last)

	return sixb.MeanI(first, last)
}

// partition slc, returns k with slc[:k] ≤ pivot ≤ slc[k:]
// swap: slc[h] < pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
//
//go:nosplit
func partOne(lsw Lesswap, l, pv, h int) int {
	// avoid unnecessary comparisons, extend ranges in balance
	for ; l < h; l, h = l+1, h-1 {

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
	}
	// classify mid element
	if l == h && h != pv && lsw(h, pv, h, h) { // 3rd=4th disables swap
		l++
	}
	return l
}

// swaps elements to get slc[lo..l] ≤ pivot ≤ slc[h..hi]
// Gap (l,h) expands until one of the intervals is fully consumed.
// swap: slc[h] < pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
//
//go:nosplit
func partTwo(lsw Lesswap, lo, l, pv, h, hi int) int {
	// avoid unnecessary comparisons, extend ranges in balance
	for {
		if lsw(h, pv, h, h) { // 3rd=4th disables swap
			for {
				if lsw(pv, l, h, l) {
					break
				}
				l--
				if l < lo {
					return h
				}
			}
		} else if lsw(pv, l, l, l) { // 3rd=4th disables swap
			for {
				h++
				if h > hi {
					return l
				}
				if lsw(h, pv, h, l) {
					break
				}
			}
		}
		l--
		h++
		if l < lo {
			return h
		}
		if h > hi {
			return l
		}
	}
}

// new-goroutine partition
//
//go:nosplit
func gPartOne(lsw Lesswap, l, pv, h int, ch chan int) {
	ch <- partOne(lsw, l, pv, h)
}

// partition slc in two goroutines, returns k with slc[:k] ≤ pivot ≤ slc[k:]
//
//go:nosplit
func partCon(lsw Lesswap, lo, hi int, ch chan int) int {

	pv := pivot(lsw, lo, hi, nsConc-1) // median-of-n pivot
	lo++
	hi--
	l, h := sixb.MeanI(lo, pv), sixb.MeanI(pv, hi)

	go gPartOne(lsw, l+1, pv, h-1, ch) // mid half range

	r := partTwo(lsw, lo, l, pv, h, hi) // left/right quarter ranges

	k := <-ch

	// only one gap is possible
	if r < pv {
		for ; lo <= r; r-- { // gap left in low range?
			if lsw(pv, r, k-1, r) {
				k--
				if k == pv { // swapped pivot when closing gap?
					pv = r // Thanks to my wife Tansu who discovered this
				}
			}
		}
	} else {
		for ; r <= hi; r++ { // gap left in high range?
			if lsw(r, pv, r, k) {
				if k == pv { // swapped pivot when closing gap?
					pv = r // It took days of agony to discover these two if's :D
				}
				k++
			}
		}
	}
	return k
}

// short range sort function, assumes MaxLenInsFC <= hi-lo < MaxLenRecFC, recursive
func short(lsw Lesswap, lo, hi int) {
start:
	fr, step, _ := minMaxSample(uint(hi+1-lo), 3)
	first := lo + int(fr)
	pv := first + int(step)
	last := pv + int(step)

	lsw(pv, first, pv, first)
	if lsw(last, pv, last, pv) {
		lsw(pv, first, pv, first) // median-of-3 pivot
	}

	// move one sorted sample to each end
	lsw(first, lo, first, lo)
	lsw(hi, last, hi, last)

	l := partOne(lsw, lo+1, pv, hi-1)
	h := l - 1
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
	// at least one insertion range, insertion inlined
isort:
	for k := l + 1; k <= h; k++ {
		for i := k; lsw(i, i-1, i, i-1); {
			i--
			if i <= l {
				break
			}
		}
	}

	if no >= MaxLenInsFC {
		goto start
	}
	if lo != l {
		l, h = lo, hi
		goto isort // two insertion ranges
	}
}

// new-goroutine sort function
//
//go:nosplit
func gLong(lsw Lesswap, lo, hi int, sv *syncVar) {
	long(lsw, lo, hi, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes hi-lo >= MaxLenRecFC, recursive
func long(lsw Lesswap, lo, hi int, sv *syncVar) {
start:
	pv := pivot(lsw, lo, hi, nsLong-1) // median-of-n pivot
	l := partOne(lsw, lo+1, pv, hi-1)
	h := l - 1
	no, n := h-lo, hi-l

	if no < n {
		n, no = no, n // [lo,hi] is the longer range
		l, lo = lo, l
	} else {
		h, hi = hi, h
	}

	// branches below are optimal for fewer total jumps
	if n < MaxLenRecFC { // at least one not-long range?

		if n >= MaxLenInsFC {
			short(lsw, l, h)
		} else {
			insertion(lsw, l, h)
		}

		if no >= MaxLenRecFC { // two not-long ranges?
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

	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	atomic.AddUint64(&sv.nGor, 1) // increase goroutine counter
	go gLong(lsw, lo, hi, sv)
	lo, hi = l, h
	goto start
}

// Sort concurrently sorts underlying collection of length n via lsw().
// Once for each non-trivial type you want to sort in a certain way, you
// can implement a custom sorting routine (for a slice for example) as:
//
//	func SortTypeAscending(slc []Type) {
//		lsw := func(i, k, r, s int) bool {
//			if slc[i].Key < slc[k].Key { // strict comparator like < or >
//				if r != s {
//					slc[r], slc[s] = slc[s], slc[r]
//				}
//				return true
//			}
//			return false
//		}
//		sorty.Sort(len(slc), lsw)
//	}
//
// [Lesswap] is a contract between users and sorty. Strict
// comparator, r!=s check, swap and returns are all necessary.
//
//go:nosplit
func Sort(n int, lsw Lesswap) {

	n-- // high index
	if n <= 2*MaxLenRecFC || MaxGor <= 1 {

		if n >= MaxLenRecFC { // single-goroutine sorting
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
		l := partCon(lsw, lo, hi, sv.done)
		h := l - 1
		no, n := h-lo, hi-l

		if no < n {
			n, no = no, n // [lo,hi] is the longer range
			l, lo = lo, l
		} else {
			h, hi = hi, h
		}

		// handle shorter range
		if n >= MaxLenRecFC {
			atomic.AddUint64(&sv.nGor, 1) // increase goroutine counter
			go gLong(lsw, l, h, &sv)

		} else if n >= MaxLenInsFC {
			short(lsw, l, h)
		} else {
			insertion(lsw, l, h)
		}

		// longer range big enough? max goroutines?
		if no <= 2*MaxLenRecFC || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	long(lsw, lo, hi, &sv) // we know hi-lo >= MaxLenRecFC

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
