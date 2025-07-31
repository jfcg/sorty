/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"cmp"

	"github.com/jfcg/sixb/v2"
)

// isSortedO returns 0 if slc is sorted in ascending order, otherwise
// it returns i > 0 with slc[i] < slc[i-1] or either is a NaN, inlined
func isSortedO[S ~[]T, T cmp.Ordered](slc S) int {
	for i := len(slc) - 1; i > 0; i-- {
		if !(slc[i] >= slc[i-1]) {
			return i
		}
	}
	return 0
}

// insertion sort, inlined
func insertionO[S ~[]T, T cmp.Ordered](slc S) {

	for h := 1; h < len(slc); h++ {
		l, val := h, slc[h]
		var pre T
		goto start
	loop:
		slc[l] = pre
		l--
		if l == 0 {
			goto last
		}
	start:
		pre = slc[l-1]
		if val < pre {
			goto loop
		}
		if l == h {
			continue
		}
	last:
		slc[l] = val
	}
}

// pivotO selects n equidistant samples from slc that minimizes max distance
// to non-selected members, then sorts the samples and returns the middle two.
// Assumes nsConc ≥ n ≥ 5, len(slc) ≥ 2n.
//
//go:nosplit
func pivotO[S ~[]T, T cmp.Ordered](slc S, n uint) (T, T) {

	first, step, last := minMaxSample(uint(len(slc)), n)

	var sample [nsConc]T
	a, b := slc[first], slc[last]
	if b < a {
		a, b = b, a
	}
	sample[0], sample[n-1] = a, b

	for i := n - 2; i > 0; i-- {
		last -= step
		sample[i] = slc[last]
	}
	insertionO(sample[:n]) // sort n samples

	n >>= 1 // return middle two samples
	return sample[n-1], sample[n]
}

// partition slc, returns k with slc[:k] ≤ pivot ≤ slc[k:]
// swap: slc[h] < pv ≤ slc[l]
// swap: slc[h] ≤ pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
//
//go:nosplit
func partOneO[S ~[]T, T cmp.Ordered](slc S, pv T) int {
	l, h := 0, len(slc)-1
	goto start
second:
	for {
		h--
		if h <= l {
			return l
		}
		if slc[h] <= pv {
			break
		}
	}
swap:
	slc[l], slc[h] = slc[h], slc[l]
next:
	l++
	h--
start:
	if h <= l {
		goto last
	}

	if pv <= slc[h] { // avoid unnecessary comparisons
		if pv < slc[l] { // extend ranges in balance
			goto second
		}
		goto next
	}
	for {
		if pv <= slc[l] {
			goto swap
		}
		l++
		if h <= l {
			return l + 1
		}
	}
last:
	if l == h && slc[h] < pv { // classify mid element
		l++
	}
	return l
}

// swaps elements to get slc[:l] ≤ pivot ≤ slc[h:]
// Gap (l,h) expands until one of the intervals is fully consumed.
// swap: slc[h] < pv ≤ slc[l]
// swap: slc[h] ≤ pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
//
//go:nosplit
func partTwoO[S ~[]T, T cmp.Ordered](slc S, l, h int, pv T) int {
	l--
	if h <= l {
		return -1 // will not run
	}
	goto start
second:
	for {
		h++
		if h >= len(slc) {
			return l
		}
		if slc[h] <= pv {
			break
		}
	}
swap:
	slc[l], slc[h] = slc[h], slc[l]
next:
	l--
	h++
start:
	if l < 0 {
		return h
	}
	if h >= len(slc) {
		return l
	}

	if pv <= slc[h] { // avoid unnecessary comparisons
		if pv < slc[l] { // extend ranges in balance
			goto second
		}
		goto next
	}
	for {
		if pv <= slc[l] {
			goto swap
		}
		l--
		if l < 0 {
			return h
		}
	}
}

// new-goroutine partition
//
//go:nosplit
func gPartOneO[S ~[]T, T cmp.Ordered](slc S, pv T, ch chan int) {
	ch <- partOneO(slc, pv)
}

// partition slc in two goroutines, returns k with slc[:k] ≤ pivot ≤ slc[k:]
//
//go:nosplit
func partConO[S ~[]T, T cmp.Ordered](slc S, pv T, ch chan int) int {
	mid := len(slc) >> 1
	l, h := mid>>1, sixb.Mean(mid, len(slc))

	go gPartOneO(slc[l:h:h], pv, ch) // mid half range

	r := partTwoO(slc, l, h, pv) // left/right quarter ranges

	k := l + <-ch // convert returned index to slc

	// only one gap is possible
	if r < mid {
		for ; 0 <= r; r-- { // gap left in low range?
			if pv < slc[r] {
				k--
				slc[r], slc[k] = slc[k], slc[r]
			}
		}
	} else {
		for ; r < len(slc); r++ { // gap left in high range?
			if slc[r] < pv {
				slc[r], slc[k] = slc[k], slc[r]
				k++
			}
		}
	}
	return k
}
