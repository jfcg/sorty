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

// isSortedF4 returns 0 if slc is sorted in ascending order, otherwise it returns i > 0
// with slc[i] < slc[i-1] or either one is a NaN. NaNoption is taken into account.
func isSortedF4(slc []float32) int {
	l, h := 0, len(slc)-1
	if NaNoption == NaNlarge { // ignore NaNs at the end
		for ; l <= h; h-- {
			if x := slc[h]; x == x {
				break
			}
		}
	} else if NaNoption == NaNsmall { // ignore NaNs at the start
		for ; l <= h; l++ {
			if x := slc[l]; x == x {
				break
			}
		}
	}

	for i := h; i > l; i-- {
		if !(slc[i] >= slc[i-1]) {
			return i
		}
	}
	return 0
}

// insertion sort, inlined
func insertionF4(slc []float32) {
	for h := 1; h < len(slc); h++ {
		l, val := h, slc[h]
		var pre float32
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

// pivotF4 selects n equidistant samples from slc that minimizes max distance
// to non-selected members, then calculates median-of-n pivot from samples.
// Assumes odd n, nsConc > n ≥ 3, len(slc) ≥ 2n. Returns pivot for partitioning.
//
//go:nosplit
func pivotF4(slc []float32, n uint) float32 {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc - 1]float32
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = slc[first]
		first += step
	}
	insertionF4(sample[:n]) // sort n samples

	return sample[n>>1] // return middle sample
}

// partition slc, returns k with slc[:k] ≤ pivot ≤ slc[k:]
// swap: slc[h] < pv ≤ slc[l]
// swap: slc[h] ≤ pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
//
//go:nosplit
func partOneF4(slc []float32, pv float32) int {
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
func partTwoF4(slc []float32, l, h int, pv float32) int {
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
func gPartOneF4(ar []float32, pv float32, ch chan int) {
	ch <- partOneF4(ar, pv)
}

// partition slc in two goroutines, returns k with slc[:k] ≤ pivot ≤ slc[k:]
//
//go:nosplit
func partConF4(slc []float32, ch chan int) int {

	pv := pivotF4(slc, nsConc-1) // median-of-n pivot
	mid := len(slc) >> 1
	l, h := mid>>1, sixb.Mean(mid, len(slc))

	go gPartOneF4(slc[l:h:h], pv, ch) // mid half range

	r := partTwoF4(slc, l, h, pv) // left/right quarter ranges

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
func shortF4(ar []float32) {
start:
	first, step, last := minMaxSample(uint(len(ar)), 3)
	pv := sixb.Median3(ar[first], ar[first+step], ar[last])

	k := partOneF4(ar, pv)
	var aq []float32

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortF4(aq) // recurse on the shorter range
		goto start
	}
isort:
	insertionF4(aq) // at least one insertion range

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
func gLongF4(ar []float32, sv *syncVar) {
	longF4(ar, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec, recursive
func longF4(ar []float32, sv *syncVar) {
start:
	pv := pivotF4(ar, nsLong-1) // median-of-n pivot
	k := partOneF4(ar, pv)
	var aq []float32

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
			shortF4(aq)
		} else {
			insertionF4(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortF4(ar) // we know len(ar) > MaxLenIns
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longF4(aq, sv) // recurse on the shorter range
		goto start
	}

	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	atomic.AddUint64(&sv.nGor, 1) // increase goroutine counter
	go gLongF4(ar, sv)
	ar = aq
	goto start
}

// sortF4 concurrently sorts ar in ascending order.
//
//go:nosplit
func sortF4(ar []float32) {
	l, h := 0, len(ar)-1
	if NaNoption == NaNlarge { // move NaNs to the end
		for l <= h {
			x := ar[h]
			if x != x {
				h--
				continue
			}
			y := ar[l]
			if y != y {
				ar[l], ar[h] = x, y
				h--
			}
			l++
		}
		ar = ar[:h+1]
	} else if NaNoption == NaNsmall { // move NaNs to the start
		for l <= h {
			y := ar[l]
			if y != y {
				l++
				continue
			}
			x := ar[h]
			if x != x {
				ar[l], ar[h] = x, y
				l++
			}
			h--
		}
		ar = ar[l:]
	}

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longF4(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortF4(ar)
		} else {
			insertionF4(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := partConF4(ar, sv.done)
		var aq []float32

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
			go gLongF4(aq, &sv)

		} else if len(aq) > MaxLenIns {
			shortF4(aq)
		} else {
			insertionF4(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longF4(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
