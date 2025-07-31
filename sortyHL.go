/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"sync/atomic"

	"github.com/jfcg/sixb/v2"
)

type hasLen interface {
	~string | ~[]struct{}
}

// isSortedHL returns 0 if ar is sorted by length in ascending
// order, otherwise it returns i > 0 with len(ar[i]) < len(ar[i-1]), inlined
func isSortedHL[S ~[]T, T hasLen](ar S) int {
	for i := len(ar) - 1; i > 0; i-- {
		if len(ar[i]) < len(ar[i-1]) {
			return i
		}
	}
	return 0
}

// insertion sort, inlined
func insertionHL[S ~[]T, T hasLen](slc S) {
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
		if len(val) < len(pre) {
			goto loop
		}
		if l == h {
			continue
		}
	last:
		slc[l] = val
	}
}

// pivotHL selects n equidistant samples from slc that minimizes max distance
// to non-selected members, then calculates median-of-n pivot from samples.
// Assumes even n, nsConc ≥ n ≥ 2, len(slc) ≥ 2n. Returns pivot for partitioning.
//
//go:nosplit
func pivotHL[S ~[]T, T hasLen](slc S, n uint) int {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc]int
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = len(slc[first])
		first += step
	}
	insertionI(sample[:n]) // sort n samples

	n >>= 1 // return mean of middle two samples
	return sixb.Mean(sample[n-1], sample[n])
}

// partition slc, returns k with slc[:k] ≤ pivot ≤ slc[k:]
// swap: slc[h] < pv ≤ slc[l]
// swap: slc[h] ≤ pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
//
//go:nosplit
func partOneHL[S ~[]T, T hasLen](slc S, pv int) int {
	l, h := 0, len(slc)-1
	goto start
second:
	for {
		h--
		if h <= l {
			return l
		}
		if len(slc[h]) <= pv {
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

	if pv <= len(slc[h]) { // avoid unnecessary comparisons
		if pv < len(slc[l]) { // extend ranges in balance
			goto second
		}
		goto next
	}
	for {
		if pv <= len(slc[l]) {
			goto swap
		}
		l++
		if h <= l {
			return l + 1
		}
	}
last:
	if l == h && len(slc[h]) < pv { // classify mid element
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
func partTwoHL[S ~[]T, T hasLen](slc S, l, h int, pv int) int {
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
		if len(slc[h]) <= pv {
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

	if pv <= len(slc[h]) { // avoid unnecessary comparisons
		if pv < len(slc[l]) { // extend ranges in balance
			goto second
		}
		goto next
	}
	for {
		if pv <= len(slc[l]) {
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
func gPartOneHL[S ~[]T, T hasLen](ar S, pv int, ch chan int) {
	ch <- partOneHL(ar, pv)
}

// partition slc in two goroutines, returns k with slc[:k] ≤ pivot ≤ slc[k:]
//
//go:nosplit
func partConHL[S ~[]T, T hasLen](slc S, ch chan int) int {

	pv := pivotHL(slc, nsConc) // median-of-n pivot
	mid := len(slc) >> 1
	l, h := mid>>1, sixb.Mean(mid, len(slc))

	go gPartOneHL(slc[l:h:h], pv, ch) // mid half range

	r := partTwoHL(slc, l, h, pv) // left/right quarter ranges

	k := l + <-ch // convert returned index to slc

	// only one gap is possible
	if r < mid {
		for ; 0 <= r; r-- { // gap left in low range?
			if pv < len(slc[r]) {
				k--
				slc[r], slc[k] = slc[k], slc[r]
			}
		}
	} else {
		for ; r < len(slc); r++ { // gap left in high range?
			if len(slc[r]) < pv {
				slc[r], slc[k] = slc[k], slc[r]
				k++
			}
		}
	}
	return k
}

// short range sort function, assumes MaxLenIns < len(ar) <= MaxLenRec, recursive
func shortHL[S ~[]T, T hasLen](ar S) {
start:
	first, step := minMaxFour(uint32(len(ar)))
	pv := sixb.Median4(len(ar[first]), len(ar[first+step]),
		len(ar[first+2*step]), len(ar[first+3*step]))

	k := partOneHL(ar, pv)
	var aq S

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortHL(aq) // recurse on the shorter range
		goto start
	}
isort:
	insertionHL(aq) // at least one insertion range

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
func gLongHL[S ~[]T, T hasLen](ar S, sv *syncVar) {
	longHL(ar, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec, recursive
func longHL[S ~[]T, T hasLen](ar S, sv *syncVar) {
start:
	pv := pivotHL(ar, nsLong) // median-of-n pivot
	k := partOneHL(ar, pv)
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
			shortHL(aq)
		} else {
			insertionHL(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortHL(ar) // we know len(ar) > MaxLenIns
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longHL(aq, sv) // recurse on the shorter range
		goto start
	}

	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	atomic.AddUint64(&sv.nGor, 1) // increase goroutine counter
	go gLongHL(ar, sv)
	ar = aq
	goto start
}

// sortHL concurrently sorts ar by length in ascending order.
//
//go:nosplit
func sortHL[S ~[]T, T hasLen](ar S) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longHL(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortHL(ar)
		} else {
			insertionHL(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := partConHL(ar, sv.done)
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
			go gLongHL(aq, &sv)

		} else if len(aq) > MaxLenIns {
			shortHL(aq)
		} else {
			insertionHL(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longHL(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
