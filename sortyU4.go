/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"sync/atomic"

	"github.com/jfcg/sixb"
)

// isSortedU4 returns 0 if ar is sorted in ascending
// order, otherwise it returns i > 0 with ar[i] < ar[i-1]
func isSortedU4(ar []uint32) int {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return i
		}
	}
	return 0
}

// insertion sort
func insertionU4(slc []uint32) {
	for h := 1; h < len(slc); h++ {
		l, val := h, slc[h]
		var pre uint32
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

// pivotU4 selects n equidistant samples from slc that minimizes max distance
// to non-selected members, then calculates median-of-n pivot from samples.
// Assumes even n, nsConc ≥ n ≥ 2, len(slc) ≥ 2n. Returns pivot for partitioning.
func pivotU4(slc []uint32, n uint) uint32 {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc]uint32
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = slc[first]
		first += step
	}
	insertionU4(sample[:n]) // sort n samples

	n >>= 1 // return mean of middle two samples
	return sixb.MeanU4(sample[n-1], sample[n])
}

// partition slc, returns k with slc[:k] ≤ pivot ≤ slc[k:]
// swap: slc[h] < pv ≤ slc[l]
// swap: slc[h] ≤ pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
func partOneU4(slc []uint32, pv uint32) int {
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
func partTwoU4(slc []uint32, l, h int, pv uint32) (int, int) {
	l--
	if h <= l {
		return l, h
	}
	goto start
second:
	for {
		h++
		if h >= len(slc) {
			return l, h
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
	if l < 0 || h >= len(slc) {
		return l, h
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
			return l, h
		}
	}
}

// new-goroutine partition
func gPartOneU4(ar []uint32, pv uint32, ch chan int) {
	ch <- partOneU4(ar, pv)
}

// partition slc in two goroutines, returns k with slc[:k] ≤ pivot ≤ slc[k:]
func partConU4(slc []uint32, ch chan int) int {

	pv := pivotU4(slc, nsConc) // median-of-n pivot
	k := len(slc) >> 1
	l, h := k>>1, sixb.MeanI(k, len(slc))

	go gPartOneU4(slc[l:h:h], pv, ch) // mid half range

	k = l
	l, h = partTwoU4(slc, l, h, pv) // left/right quarter ranges

	k += <-ch // convert returned indice to slc

	// only one gap is possible
	for ; 0 <= l; l-- { // gap left in low range?
		if pv < slc[l] {
			k--
			slc[l], slc[k] = slc[k], slc[l]
		}
	}
	for ; h < len(slc); h++ { // gap left in high range?
		if slc[h] < pv {
			slc[h], slc[k] = slc[k], slc[h]
			k++
		}
	}
	return k
}

// short range sort function, assumes MaxLenIns < len(ar) <= MaxLenRec
func shortU4(ar []uint32) {
start:
	pv := pivotU4(ar, nsShort) // median-of-n pivot
	k := partOneU4(ar, pv)
	var aq []uint32

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortU4(aq) // recurse on the shorter range
		goto start
	}
psort:
	insertionU4(aq) // at least one insertion range

	if len(ar) > MaxLenIns {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto psort // two insertion ranges
	}
}

// new-goroutine sort function
func glongU4(ar []uint32, sv *syncVar) {
	longU4(ar, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec
func longU4(ar []uint32, sv *syncVar) {
start:
	pv := pivotU4(ar, nsLong) // median-of-n pivot
	k := partOneU4(ar, pv)
	var aq []uint32

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
			shortU4(aq)
		} else {
			insertionU4(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortU4(ar) // we know len(ar) > MaxLenIns
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longU4(aq, sv) // recurse on the shorter range
		goto start
	}

	if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
		panic("sorty: longU4: counter overflow")
	}
	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	go glongU4(ar, sv)
	ar = aq
	goto start
}

// sortU4 concurrently sorts ar in ascending order.
func sortU4(ar []uint32) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longU4(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortU4(ar)
		} else {
			insertionU4(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := partConU4(ar, sv.done)
		var aq []uint32

		if k < len(ar)-k {
			aq = ar[:k:k]
			ar = ar[k:] // ar is the longer range
		} else {
			aq = ar[k:]
			ar = ar[:k:k]
		}

		// handle shorter range
		if len(aq) > MaxLenRec {
			if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
				panic("sorty: sortU4: counter overflow")
			}
			go glongU4(aq, &sv)

		} else if len(aq) > MaxLenIns {
			shortU4(aq)
		} else {
			insertionU4(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longU4(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
