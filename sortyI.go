/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"sync/atomic"

	sb "github.com/jfcg/sixb/v2"
)

// isSortedI returns 0 if slc is sorted in ascending order
// otherwise it returns i > 0 with slc[i] < slc[i-1], inlined
func isSortedI[S ~[]T, T sb.Integer](slc S) int {
	for i := len(slc) - 1; i > 0; i-- {
		if slc[i] < slc[i-1] {
			return i
		}
	}
	return 0
}

// insertion sort, inlined
func insertionI[S ~[]T, T sb.Integer](slc S) {
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

// pivotI selects n equidistant samples from slc that minimizes max distance
// to non-selected members, then calculates median-of-n pivot from samples.
// Assumes even n, nsConc ≥ n ≥ 2, len(slc) ≥ 2n. Returns pivot for partitioning.
//
//go:nosplit
func pivotI[S ~[]T, T sb.Integer](slc S, n uint) T {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc]T
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = slc[first]
		first += step
	}
	insertionI(sample[:n]) // sort n samples

	n >>= 1 // return mean of middle two samples
	return sb.Mean(sample[n-1], sample[n])
}

// partition slc, returns k with slc[:k] ≤ pivot ≤ slc[k:]
// swap: slc[h] < pv ≤ slc[l]
// swap: slc[h] ≤ pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
//
//go:nosplit
func partOneI[S ~[]T, T sb.Integer](slc S, pv T) int {
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
func partTwoI[S ~[]T, T sb.Integer](slc S, l, h int, pv T) int {
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
func gPartOneI[S ~[]T, T sb.Integer](slc S, pv T, ch chan int) {
	ch <- partOneI(slc, pv)
}

// partition slc in two goroutines, returns k with slc[:k] ≤ pivot ≤ slc[k:]
//
//go:nosplit
func partConI[S ~[]T, T sb.Integer](slc S, ch chan int) int {

	pv := pivotI(slc, nsConc) // median-of-n pivot
	mid := len(slc) >> 1
	l, h := mid>>1, sb.Mean(mid, len(slc))

	go gPartOneI(slc[l:h:h], pv, ch) // mid half range

	r := partTwoI(slc, l, h, pv) // left/right quarter ranges

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

// short range sort function, assumes MaxLenIns < len(ar) <= MaxLenRec, recursive
func shortI[S ~[]T, T sb.Integer](ar S) {
start:
	first, step := minMaxFour(uint32(len(ar)))
	pv := sb.Median4(ar[first], ar[first+step], ar[first+2*step], ar[first+3*step])

	k := partOneI(ar, pv)
	var aq S

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortI(aq) // recurse on the shorter range
		goto start
	}
isort:
	insertionI(aq) // at least one insertion range

	if len(ar) > MaxLenIns {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto isort // two insertion ranges
	}
}

// new-goroutine sort function
//
//go:nosplit
func gLongI[S ~[]T, T sb.Integer](ar S, sv *syncVar) {
	longI(ar, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec, recursive
func longI[S ~[]T, T sb.Integer](ar S, sv *syncVar) {
start:
	pv := pivotI(ar, nsLong) // median-of-n pivot
	k := partOneI(ar, pv)
	var aq S

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	// branches below are optimal for fewer total jumps
	if len(aq) <= MaxLenRec { // at least one not-long range?

		if len(aq) > MaxLenIns {
			shortI(aq)
		} else {
			insertionI(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortI(ar) // we know len(ar) > MaxLenIns
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longI(aq, sv) // recurse on the shorter range
		goto start
	}

	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	atomic.AddUint64(&sv.nGor, 1) // increase goroutine counter
	go gLongI(ar, sv)
	ar = aq
	goto start
}

// sortI concurrently sorts ar in ascending order.
//
//go:nosplit
func sortI[S ~[]T, T sb.Integer](ar S) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longI(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortI(ar)
		} else {
			insertionI(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := partConI(ar, sv.done)
		var aq S

		if k < len(ar)-k {
			aq = ar[:k:k]
			ar = ar[k:] // ar is the longer range
		} else {
			aq = ar[k:]
			ar = ar[:k:k]
		}

		// handle shorter range
		if len(aq) > MaxLenRec {
			atomic.AddUint64(&sv.nGor, 1) // increase goroutine counter
			go gLongI(aq, &sv)

		} else if len(aq) > MaxLenIns {
			shortI(aq)
		} else {
			insertionI(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longI(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
