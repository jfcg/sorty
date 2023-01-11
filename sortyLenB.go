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

// isSortedLenB returns 0 if ar is sorted by length in ascending
// order, otherwise it returns i > 0 with len(ar[i]) < len(ar[i-1])
func isSortedLenB(ar [][]byte) int {
	for i := len(ar) - 1; i > 0; i-- {
		if len(ar[i]) < len(ar[i-1]) {
			return i
		}
	}
	return 0
}

// insertion sort
func insertionLenB(slc [][]byte) {
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

// pivotLenB selects n equidistant samples from slc that minimizes max distance
// to non-selected members, then calculates median-of-n pivot from samples.
// Assumes even n, nsConc ≥ n ≥ 2, len(slc) ≥ 2n. Returns pivot for partitioning.
func pivotLenB(slc [][]byte, n uint) int {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc]int
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = len(slc[first])
		first += step
	}
	insertionI(sample[:n]) // sort n samples

	n >>= 1 // return mean of middle two samples
	return sixb.MeanI(sample[n-1], sample[n])
}

// partition ar into <= and >= pivot, assumes len(ar) >= 2
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func partition1LenB(ar [][]byte, pv int) int {
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
func partition2LenB(ar [][]byte, a, b int, pv int) (int, int) {
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
func gpart1LenB(ar [][]byte, pv int, ch chan int) {
	ch <- partition1LenB(ar, pv)
}

// concurrent dual partitioning of ar
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func cdualparLenB(ar [][]byte, ch chan int) int {

	pv := pivotLenB(ar, nsConc) // median-of-n pivot
	k := len(ar) >> 1
	a, b := k>>1, sixb.MeanI(k, len(ar))

	go gpart1LenB(ar[a:b:b], pv, ch) // mid half range

	k = a
	a, b = partition2LenB(ar, a, b, pv) // left/right quarter ranges

	k += <-ch // convert returned indice to ar

	// only one gap is possible
	for ; 0 <= a; a-- { // gap left in low range?
		if pv < len(ar[a]) {
			k--
			ar[a], ar[k] = ar[k], ar[a]
		}
	}
	for ; b < len(ar); b++ { // gap left in high range?
		if len(ar[b]) < pv {
			ar[b], ar[k] = ar[k], ar[b]
			k++
		}
	}
	return k
}

// short range sort function, assumes MaxLenIns < len(ar) <= MaxLenRec
func shortLenB(ar [][]byte) {
start:
	pv := pivotLenB(ar, nsShort) // median-of-n pivot
	k := partition1LenB(ar, pv)
	var aq [][]byte

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortLenB(aq) // recurse on the shorter range
		goto start
	}
psort:
	insertionLenB(aq) // at least one insertion range

	if len(ar) > MaxLenIns {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto psort // two insertion ranges
	}
}

// new-goroutine sort function
func glongLenB(ar [][]byte, sv *syncVar) {
	longLenB(ar, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec
func longLenB(ar [][]byte, sv *syncVar) {
start:
	pv := pivotLenB(ar, nsLong) // median-of-n pivot
	k := partition1LenB(ar, pv)
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

		if len(aq) > MaxLenIns {
			shortLenB(aq)
		} else {
			insertionLenB(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortLenB(ar) // we know len(ar) > MaxLenIns
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longLenB(aq, sv) // recurse on the shorter range
		goto start
	}

	if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
		panic("sorty: longLenB: counter overflow")
	}
	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	go glongLenB(ar, sv)
	ar = aq
	goto start
}

// sortLenB concurrently sorts ar by length in ascending order.
func sortLenB(ar [][]byte) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longLenB(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortLenB(ar)
		} else {
			insertionLenB(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := cdualparLenB(ar, sv.done)
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
				panic("sorty: sortLenB: counter overflow")
			}
			go glongLenB(aq, &sv)

		} else if len(aq) > MaxLenIns {
			shortLenB(aq)
		} else {
			insertionLenB(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longLenB(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
