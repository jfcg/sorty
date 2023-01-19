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

// isSortedB returns 0 if ar is sorted in ascending lexicographic
// order, otherwise it returns i > 0 with sixb.BtoS(ar[i]) < sixb.BtoS(ar[i-1])
func isSortedB(ar [][]byte) int {
	for i := len(ar) - 1; i > 0; i-- {
		if sixb.BtoS(ar[i]) < sixb.BtoS(ar[i-1]) {
			return i
		}
	}
	return 0
}

// insertion sort
func insertionB(slc [][]byte) {
	for h := 1; h < len(slc); h++ {
		l, val := h, slc[h]
		var pre []byte
		goto start
	loop:
		slc[l] = pre
		l--
		if l == 0 {
			goto last
		}
	start:
		pre = slc[l-1]
		if sixb.BtoS(val) < sixb.BtoS(pre) {
			goto loop
		}
		if l == h {
			continue
		}
	last:
		slc[l] = val
	}
}

// pivotB selects n equidistant samples from slc that minimizes max distance
// to non-selected members, then calculates median-of-n pivot from samples.
// Assumes odd n, nsConc > n ≥ 3, len(slc) ≥ 2n. Returns pivot for partitioning.
func pivotB(slc [][]byte, n uint) string {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc - 1]string
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = sixb.BtoS(slc[first])
		first += step
	}
	insertionS(sample[:n]) // sort n samples

	return sample[n>>1] // return middle sample
}

// partition slc, returns k with slc[:k] ≤ pivot ≤ slc[k:]
// swap: slc[h] < pv ≤ slc[l]
// swap: slc[h] ≤ pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
func partOneB(slc [][]byte, pv string) int {
	l, h := 0, len(slc)-1
	goto start
second:
	for {
		h--
		if h <= l {
			return l
		}
		if sixb.BtoS(slc[h]) <= pv {
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

	if pv <= sixb.BtoS(slc[h]) { // avoid unnecessary comparisons
		if pv < sixb.BtoS(slc[l]) { // extend ranges in balance
			goto second
		}
		goto next
	}
	for {
		if pv <= sixb.BtoS(slc[l]) {
			goto swap
		}
		l++
		if h <= l {
			return l + 1
		}
	}
last:
	if l == h && sixb.BtoS(slc[h]) < pv { // classify mid element
		l++
	}
	return l
}

// swaps elements to get slc[:l] ≤ pivot ≤ slc[h:]
// Gap (l,h) expands until one of the intervals is fully consumed.
// swap: slc[h] < pv ≤ slc[l]
// swap: slc[h] ≤ pv < slc[l]
// next: slc[l] ≤ pv ≤ slc[h]
func partTwoB(slc [][]byte, l, h int, pv string) (int, int) {
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
		if sixb.BtoS(slc[h]) <= pv {
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

	if pv <= sixb.BtoS(slc[h]) { // avoid unnecessary comparisons
		if pv < sixb.BtoS(slc[l]) { // extend ranges in balance
			goto second
		}
		goto next
	}
	for {
		if pv <= sixb.BtoS(slc[l]) {
			goto swap
		}
		l--
		if l < 0 {
			return l, h
		}
	}
}

// new-goroutine partition
func gPartOneB(ar [][]byte, pv string, ch chan int) {
	ch <- partOneB(ar, pv)
}

// partition slc in two goroutines, returns k with slc[:k] ≤ pivot ≤ slc[k:]
func partConB(slc [][]byte, ch chan int) int {

	pv := pivotB(slc, nsConc-1) // median-of-n pivot
	k := len(slc) >> 1
	l, h := k>>1, sixb.MeanI(k, len(slc))

	go gPartOneB(slc[l:h:h], pv, ch) // mid half range

	k = l
	l, h = partTwoB(slc, l, h, pv) // left/right quarter ranges

	k += <-ch // convert returned indice to slc

	// only one gap is possible
	for ; 0 <= l; l-- { // gap left in low range?
		if pv < sixb.BtoS(slc[l]) {
			k--
			slc[l], slc[k] = slc[k], slc[l]
		}
	}
	for ; h < len(slc); h++ { // gap left in high range?
		if sixb.BtoS(slc[h]) < pv {
			slc[h], slc[k] = slc[k], slc[h]
			k++
		}
	}
	return k
}

// short range sort function, assumes MaxLenInsFC < len(ar) <= MaxLenRec
func shortB(ar [][]byte) {
start:
	first, step, last := minMaxSample(uint(len(ar)), 3)
	f, pv, l := sixb.BtoS(ar[first]), sixb.BtoS(ar[first+step]), sixb.BtoS(ar[last])

	if pv < f {
		pv, f = f, pv
	}
	if l < pv {
		if l < f {
			pv = f
		} else {
			pv = l // median-of-3 pivot
		}
	}

	k := partOneB(ar, pv)
	var aq [][]byte

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenInsFC {
		shortB(aq) // recurse on the shorter range
		goto start
	}
isort:
	insertionB(aq) // at least one insertion range

	if len(ar) > MaxLenInsFC {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto isort // two insertion ranges
	}
}

// new-goroutine sort function
func glongB(ar [][]byte, sv *syncVar) {
	longB(ar, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec
func longB(ar [][]byte, sv *syncVar) {
start:
	pv := pivotB(ar, nsLong-1) // median-of-n pivot
	k := partOneB(ar, pv)
	var aq [][]byte

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	// branches below are optimal for fewer total jumps
	if len(aq) <= MaxLenRec { // at least one not-long range?

		if len(aq) > MaxLenInsFC {
			shortB(aq)
		} else {
			insertionB(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortB(ar) // we know len(ar) > MaxLenInsFC
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longB(aq, sv) // recurse on the shorter range
		goto start
	}

	if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
		panic("sorty: longB: counter overflow")
	}
	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	go glongB(ar, sv)
	ar = aq
	goto start
}

// sortB concurrently sorts ar in ascending lexicographic order.
func sortB(ar [][]byte) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longB(ar, nil)
		} else if len(ar) > MaxLenInsFC {
			shortB(ar)
		} else {
			insertionB(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := partConB(ar, sv.done)
		var aq [][]byte

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
				panic("sorty: sortB: counter overflow")
			}
			go glongB(aq, &sv)

		} else if len(aq) > MaxLenInsFC {
			shortB(aq)
		} else {
			insertionB(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longB(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
