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

// short range sort function, assumes MaxLenInsFC < len(ar) <= MaxLenRecFC, recursive
func shortS(ar []string) {
start:
	first, step, last := minMaxSample(uint(len(ar)), 3)
	pv := sixb.Median3(ar[first], ar[first+step], ar[last])

	k := partOneO(ar, pv)
	var aq []string

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenInsFC {
		shortS(aq) // recurse on the shorter range
		goto start
	}
isort:
	insertionO(aq) // at least one insertion range

	if len(ar) > MaxLenInsFC {
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
func gLongS(ar []string, sv *syncVar) {
	longS(ar, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRecFC, recursive
func longS(ar []string, sv *syncVar) {
start:
	_, pv := pivotO(ar, nsLong-1) // median-of-n pivot
	k := partOneO(ar, pv)
	var aq []string

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	// branches below are optimal for fewer total jumps
	if len(aq) <= MaxLenRecFC { // at least one not-long range?

		if len(aq) > MaxLenInsFC {
			shortS(aq)
		} else {
			insertionO(aq)
		}

		if len(ar) > MaxLenRecFC { // two not-long ranges?
			goto start
		}
		shortS(ar) // we know len(ar) > MaxLenInsFC
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longS(aq, sv) // recurse on the shorter range
		goto start
	}

	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	atomic.AddUint64(&sv.nGor, 1) // increase goroutine counter
	go gLongS(ar, sv)
	ar = aq
	goto start
}

// sortS concurrently sorts ar in ascending lexicographic order.
func sortS(ar []string) {

	if len(ar) < 2*(MaxLenRecFC+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRecFC { // single-goroutine sorting
			longS(ar, nil)
		} else if len(ar) > MaxLenInsFC {
			shortS(ar)
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
		var aq []string

		if k < len(ar)-k {
			aq = ar[:k:k]
			ar = ar[k:] // ar is the longer range
		} else {
			aq = ar[k:]
			ar = ar[:k:k]
		}

		// handle shorter range
		if len(aq) > MaxLenRecFC {
			atomic.AddUint64(&sv.nGor, 1) // increase goroutine counter
			go gLongS(aq, &sv)

		} else if len(aq) > MaxLenInsFC {
			shortS(aq)
		} else {
			insertionO(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRecFC+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longS(ar, &sv) // we know len(ar) > MaxLenRecFC

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
