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

// pivotI selects n equidistant samples from slc that minimizes max distance
// to non-selected members, then sorts the samples and returns the median.
// Assumes even n with nsConc ≥ n ≥ 5, len(slc) ≥ 2n.
//
//go:nosplit
func pivotI[S ~[]T, T sb.Integer](slc S, n uint) T {
	a, b := pivotO(slc, n)
	return sb.Mean(a, b)
}

// short range sort function, assumes MaxLenIns < len(ar) <= MaxLenRec, recursive
func shortI[S ~[]T, T sb.Integer](ar S) {
start:
	first, step := minMaxFour(uint32(len(ar)))
	pv := sb.Median4(ar[first], ar[first+step], ar[first+2*step], ar[first+3*step])

	k := partOneO(ar, pv)
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
	insertionO(aq) // at least one insertion range

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
	k := partOneO(ar, pv)
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
			insertionO(aq)
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
			insertionO(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		pv := pivotI(ar, nsConc) // median-of-n pivot
		k := partConO(ar, pv, sv.done)
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
			insertionO(aq)
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
