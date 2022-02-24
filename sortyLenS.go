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

// isSortedLenS returns 0 if ar is sorted by length in ascending
// order, otherwise it returns i > 0 with len(ar[i]) < len(ar[i-1])
func isSortedLenS(ar []string) int {
	for i := len(ar) - 1; i > 0; i-- {
		if len(ar[i]) < len(ar[i-1]) {
			return i
		}
	}
	return 0
}

// pre-sort
func presortLenS(ar []string) {
	l, h := 0, (MaxLenIns+1)/3
	for h < len(ar) {
		if len(ar[h]) < len(ar[l]) {
			ar[h], ar[l] = ar[l], ar[h]
		}
		l++
		h++
	}
}

// insertion sort
func insertionLenS(ar []string) {
	for h := 0; h < len(ar)-1; {
		l := h
		h++
		v := ar[h]
		if len(v) < len(ar[l]) {
			for {
				ar[l+1] = ar[l]
				l--
				if l < 0 || len(v) >= len(ar[l]) {
					break
				}
			}
			ar[l+1] = v
		}
	}
}

// pre+insertion sort
func pinsertLenS(ar []string) {
	presortLenS(ar)
	insertionLenS(ar)
}

// pivotLenS selects 2n equidistant samples from ar that minimizes max distance to any
// non-selected member, calculates median-of-2n pivot from samples. ensures lo/hi ranges
// have at least n elements by moving sorted samples to n positions at lo/hi ends.
// assumes 5 > n > 0, len(ar) > 4n. returns remaining slice,pivot for partitioning.
func pivotLenS(ar []string, n int) ([]string, int) {

	// sample step, first/last sample positions
	d, s, h, l := minmaxSample(len(ar), n)

	var sample [8]string
	for i, k := d, h; i >= 0; i-- {
		sample[i] = ar[k]
		k -= s
	}
	insertionLenS(sample[:d+1]) // sort 2n samples

	lo, hi := 0, len(ar)
	for { // move sorted samples to lo/hi ends
		hi--
		ar[h] = ar[hi]
		ar[hi] = sample[d]
		ar[l] = ar[lo]
		ar[lo] = sample[lo]
		l += s
		h -= s
		lo++
		d--
		if d < lo {
			break
		}
	}
	return ar[lo:hi:hi], sixb.MeanI(len(sample[n-1]), len(sample[n]))
}

// partition ar into <= and >= pivot, assumes len(ar) >= 2
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func partition1LenS(ar []string, pv int) int {
	l, h := 0, len(ar)-1
	for l < h {
		if len(ar[h]) < pv { // avoid unnecessary comparisons
			for {
				if pv < len(ar[l]) {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if pv < len(ar[l]) { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l
				}
				if len(ar[h]) < pv {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
			}
		}
		l++
		h--
	}
	if l == h && len(ar[h]) < pv { // classify mid element
		l++
	}
	return l
}

// rearrange ar[:a] and ar[b:] into <= and >= pivot, assumes 0 < a < b < len(ar)
// gap (a,b) expands until one of the intervals is fully consumed
func partition2LenS(ar []string, a, b int, pv int) (int, int) {
	a--
	for a >= 0 && b < len(ar) {
		if len(ar[b]) < pv { // avoid unnecessary comparisons
			for {
				if pv < len(ar[a]) {
					ar[a], ar[b] = ar[b], ar[a]
					break
				}
				a--
				if a < 0 {
					return a, b
				}
			}
		} else if pv < len(ar[a]) { // extend ranges in balance
			for {
				b++
				if b >= len(ar) {
					return a, b
				}
				if len(ar[b]) < pv {
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
func gpart1LenS(ar []string, pv int, ch chan int) {
	ch <- partition1LenS(ar, pv)
}

// concurrent dual partitioning of ar
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func cdualparLenS(ar []string, ch chan int) int {

	aq, pv := pivotLenS(ar, 4) // median-of-8 pivot
	k := len(aq) >> 1
	a, b := k>>1, sixb.MeanI(k, len(aq))

	go gpart1LenS(aq[a:b:b], pv, ch) // mid half range

	k = a
	a, b = partition2LenS(aq, a, b, pv) // left/right quarter ranges

	k += <-ch // convert returned indice to aq

	// only one gap is possible
	for ; 0 <= a; a-- { // gap left in low range?
		if pv < len(aq[a]) {
			k--
			aq[a], aq[k] = aq[k], aq[a]
		}
	}
	for ; b < len(aq); b++ { // gap left in high range?
		if len(aq[b]) < pv {
			aq[b], aq[k] = aq[k], aq[b]
			k++
		}
	}
	return k + 4 // convert k indice to ar
}

// short range sort function, assumes MaxLenIns < len(ar) <= MaxLenRec
func shortLenS(ar []string) {
start:
	aq, pv := pivotLenS(ar, 2) // median-of-4 pivot
	k := partition1LenS(aq, pv)

	k += 2 // convert k indice from aq to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortLenS(aq) // recurse on the shorter range
		goto start
	}
psort:
	presortLenS(aq)
	insertionLenS(aq) // at least one insertion range

	if len(ar) > MaxLenIns {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto psort // two insertion ranges
	}
}

// new-goroutine sort function
func glongLenS(ar []string, sv *syncVar) {
	longLenS(ar, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec
func longLenS(ar []string, sv *syncVar) {
start:
	aq, pv := pivotLenS(ar, 3) // median-of-6 pivot
	k := partition1LenS(aq, pv)

	k += 3 // convert k indice from aq to ar

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
			shortLenS(aq)
		} else {
			pinsertLenS(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortLenS(ar) // we know len(ar) > MaxLenIns
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longLenS(aq, sv) // recurse on the shorter range
		goto start
	}

	if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
		panic("sorty: longLenS: counter overflow")
	}
	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	go glongLenS(ar, sv)
	ar = aq
	goto start
}

// sortLenS concurrently sorts ar by length in ascending order.
func sortLenS(ar []string) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longLenS(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortLenS(ar)
		} else {
			pinsertLenS(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := cdualparLenS(ar, sv.done)
		var aq []string

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
				panic("sorty: sortLenS: counter overflow")
			}
			go glongLenS(aq, &sv)

		} else if len(aq) > MaxLenIns {
			shortLenS(aq)
		} else {
			pinsertLenS(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longLenS(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
