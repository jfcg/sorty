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

// isSortedF returns 0 if slc is sorted in ascending order, otherwise it returns i > 0
// with slc[i] < slc[i-1] or either one is a NaN. NaNoption is taken into account.
func isSortedF[S ~[]T, T sb.Float](slc S) int {
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
	return isSortedO(slc[l : h+1])
}

// short range sort function, assumes MaxLenIns < len(ar) <= MaxLenRec, recursive
func shortF[S ~[]T, T sb.Float](ar S) {
start:
	first, step, last := minMaxSample(uint(len(ar)), 3)
	pv := sb.Median3(ar[first], ar[first+step], ar[last])

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
		shortF(aq) // recurse on the shorter range
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
func gLongF[S ~[]T, T sb.Float](ar S, sv *syncVar) {
	longF(ar, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec, recursive
func longF[S ~[]T, T sb.Float](ar S, sv *syncVar) {
start:
	_, pv := pivotO(ar, nsLong-1) // median-of-n pivot
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
			shortF(aq)
		} else {
			insertionO(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortF(ar) // we know len(ar) > MaxLenIns
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longF(aq, sv) // recurse on the shorter range
		goto start
	}

	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	atomic.AddUint64(&sv.nGor, 1) // increase goroutine counter
	go gLongF(ar, sv)
	ar = aq
	goto start
}

// sortF concurrently sorts ar in ascending order.
//
//go:nosplit
func sortF[S ~[]T, T sb.Float](ar S) {
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
			longF(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortF(ar)
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
		_, pv := pivotO(ar, nsConc-1) // median-of-n pivot
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
			go gLongF(aq, &sv)

		} else if len(aq) > MaxLenIns {
			shortF(aq)
		} else {
			insertionO(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longF(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
