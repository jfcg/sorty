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

// isSortedI8 returns 0 if ar is sorted in ascending
// order, otherwise it returns i > 0 with ar[i] < ar[i-1]
func isSortedI8(ar []int64) int {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return i
		}
	}
	return 0
}

// insertion sort
func insertionI8(slc []int64) {
	for h := 1; h < len(slc); h++ {
		l, val := h, slc[h]
		var pre int64
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

// pivotI8 selects n equidistant samples from slc that minimizes max distance
// to non-selected members, then calculates median-of-n pivot from samples.
// Assumes even n, nsConc ≥ n ≥ 2, len(slc) ≥ 2n. Returns pivot for partitioning.
func pivotI8(slc []int64, n uint) int64 {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc]int64
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = slc[first]
		first += step
	}
	insertionI8(sample[:n]) // sort n samples

	n >>= 1 // return mean of middle two samples
	return sixb.MeanI8(sample[n-1], sample[n])
}

// partition ar into <= and >= pivot, assumes len(ar) >= 2
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func partition1I8(ar []int64, pv int64) int {
	l, h := 0, len(ar)-1
	for l < h {
		if ar[h] < pv { // avoid unnecessary comparisons
			for {
				if pv < ar[l] {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if pv < ar[l] { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l
				}
				if ar[h] < pv {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
			}
		}
		l++
		h--
	}
	if l == h && ar[h] < pv { // classify mid element
		l++
	}
	return l
}

// rearrange ar[:a] and ar[b:] into <= and >= pivot, assumes 0 < a < b < len(ar)
// gap (a,b) expands until one of the intervals is fully consumed
func partition2I8(ar []int64, a, b int, pv int64) (int, int) {
	a--
	for a >= 0 && b < len(ar) {
		if ar[b] < pv { // avoid unnecessary comparisons
			for {
				if pv < ar[a] {
					ar[a], ar[b] = ar[b], ar[a]
					break
				}
				a--
				if a < 0 {
					return a, b
				}
			}
		} else if pv < ar[a] { // extend ranges in balance
			for {
				b++
				if b >= len(ar) {
					return a, b
				}
				if ar[b] < pv {
					ar[a], ar[b] = ar[b], ar[a]
					break
				}
			}
		}
		a--
		b++
	}
	return a, b
}

// new-goroutine partition
func gpart1I8(ar []int64, pv int64, ch chan int) {
	ch <- partition1I8(ar, pv)
}

// concurrent dual partitioning of ar
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func cdualparI8(ar []int64, ch chan int) int {

	pv := pivotI8(ar, nsConc) // median-of-n pivot
	k := len(ar) >> 1
	a, b := k>>1, sixb.MeanI(k, len(ar))

	go gpart1I8(ar[a:b:b], pv, ch) // mid half range

	k = a
	a, b = partition2I8(ar, a, b, pv) // left/right quarter ranges

	k += <-ch // convert returned indice to ar

	// only one gap is possible
	for ; 0 <= a; a-- { // gap left in low range?
		if pv < ar[a] {
			k--
			ar[a], ar[k] = ar[k], ar[a]
		}
	}
	for ; b < len(ar); b++ { // gap left in high range?
		if ar[b] < pv {
			ar[b], ar[k] = ar[k], ar[b]
			k++
		}
	}
	return k
}

// short range sort function, assumes MaxLenIns < len(ar) <= MaxLenRec
func shortI8(ar []int64) {
start:
	pv := pivotI8(ar, nsShort) // median-of-n pivot
	k := partition1I8(ar, pv)
	var aq []int64

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortI8(aq) // recurse on the shorter range
		goto start
	}
psort:
	insertionI8(aq) // at least one insertion range

	if len(ar) > MaxLenIns {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto psort // two insertion ranges
	}
}

// new-goroutine sort function
func glongI8(ar []int64, sv *syncVar) {
	longI8(ar, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec
func longI8(ar []int64, sv *syncVar) {
start:
	pv := pivotI8(ar, nsLong) // median-of-n pivot
	k := partition1I8(ar, pv)
	var aq []int64

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
			shortI8(aq)
		} else {
			insertionI8(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortI8(ar) // we know len(ar) > MaxLenIns
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longI8(aq, sv) // recurse on the shorter range
		goto start
	}

	if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
		panic("sorty: longI8: counter overflow")
	}
	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	go glongI8(ar, sv)
	ar = aq
	goto start
}

// sortI8 concurrently sorts ar in ascending order.
func sortI8(ar []int64) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longI8(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortI8(ar)
		} else {
			insertionI8(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := cdualparI8(ar, sv.done)
		var aq []int64

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
				panic("sorty: sortI8: counter overflow")
			}
			go glongI8(aq, &sv)

		} else if len(aq) > MaxLenIns {
			shortI8(aq)
		} else {
			insertionI8(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longI8(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
