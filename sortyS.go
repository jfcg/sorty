/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"sync/atomic"

	"github.com/jfcg/sixb"
)

// isSortedS returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func isSortedS(ar []string) int {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return i
		}
	}
	return 0
}

// insertion sort, assumes len(ar) >= 2
func insertionS(ar []string) {
	hi := len(ar) - 1
	for l, h := (hi-3)>>1, hi; l >= 0; {
		if ar[h] < ar[l] {
			ar[l], ar[h] = ar[h], ar[l]
		}
		l--
		h--
	}
	for h := 0; ; {
		l := h
		h++
		v := ar[h]
		if v < ar[l] {
			for {
				ar[l+1] = ar[l]
				l--
				if l < 0 || v >= ar[l] {
					break
				}
			}
			ar[l+1] = v
		}
		if h >= hi {
			break
		}
	}
}

// pivotS selects 2n equidistant samples from ar that minimizes max distance to any
// non-selected member, calculates median-of-2n pivot from samples. ensures lo/hi ranges
// have at least n elements by moving sorted samples to n positions at lo/hi ends.
// assumes 5 > n > 0, len(ar) > 4n. returns remaining slice,pivot for partitioning.
func pivotS(ar []string, n int) ([]string, string) {

	d := 2 * n
	s := len(ar) / d // sample step > 1
	d--
	h := d * s
	l := (len(ar) - h) >> 1
	if l >= n && l > (s+1)>>1 {
		s++
		h += d
		l -= n
	}
	h += l // first/last sample positions

	var sample [8]string
	for i, k := d, h; i >= 0; i-- {
		sample[i] = ar[k]
		k -= s
	}
	insertionS(sample[:2*n]) // sort 2n samples

	i, lo, hi := 0, 0, len(ar)

	// move sorted samples to lo/hi ends
	for {
		hi--
		ar[h] = ar[hi]
		ar[hi] = sample[d]
		ar[l] = ar[lo]
		ar[lo] = sample[i]
		i++
		d--
		l += s
		h -= s
		lo++
		if d < i {
			break
		}
	}
	return ar[lo:hi:hi], sixb.MeanS(sample[n-1], sample[n])
}

// partition ar into <= and >= pivot, assumes len(ar) >= 2
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func partition1S(ar []string, pv string) int {
	l, h := 0, len(ar)-1
	for {
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
		if l >= h {
			break
		}
	}
	if l == h && ar[h] < pv { // classify mid element
		l++
	}
	return l
}

// rearrange ar[:a] and ar[b:] into <= and >= pivot, assumes 0 < a < b < len(ar)
// gap (a,b) expands until one of the intervals is fully consumed
func partition2S(ar []string, a, b int, pv string) (int, int) {
	a--
	for {
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
		if a < 0 || b >= len(ar) {
			return a, b
		}
	}
}

// new-goroutine partition
func gpart1S(ar []string, pv string, ch chan int) {
	ch <- partition1S(ar, pv)
}

// concurrent dual partitioning of ar
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func cdualparS(ar []string, ch chan int) int {

	aq, pv := pivotS(ar, 4) // median-of-9
	k := len(aq) >> 1
	a, b := k>>1, mid(k, len(aq))

	go gpart1S(aq[a:b:b], pv, ch) // mid half range

	t := a
	a, b = partition2S(aq, a, b, pv) // left/right quarter ranges
	k = <-ch
	k += t // convert k indice to aq

	// only one gap is possible
	for ; 0 <= a; a-- { // gap left in low range?
		if pv < aq[a] {
			k--
			aq[a], aq[k] = aq[k], aq[a]
		}
	}
	for ; b < len(aq); b++ { // gap left in high range?
		if aq[b] < pv {
			aq[b], aq[k] = aq[k], aq[b]
			k++
		}
	}
	return k + 4 // convert k indice to ar
}

// short range sort function, assumes MaxLenInsFC < len(ar) <= MaxLenRec
func shortS(ar []string) {
start:
	aq, pv := pivotS(ar, 2)
	k := partition1S(aq, pv) // median-of-5 partitioning

	k += 2 // convert k indice from aq to ar

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
	insertionS(aq) // at least one insertion range

	if len(ar) > MaxLenInsFC {
		goto start
	}
	insertionS(ar) // two insertion ranges
}

// long range sort function (single goroutine), assumes len(ar) > MaxLenRec
func slongS(ar []string) {
start:
	aq, pv := pivotS(ar, 3)
	k := partition1S(aq, pv) // median-of-7 partitioning

	k += 3 // convert k indice from aq to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenRec { // at least one not-long range?
		slongS(aq) // recurse on the shorter range
		goto start
	}

	if len(aq) > MaxLenInsFC {
		shortS(aq)
	} else {
		insertionS(aq)
	}

	if len(ar) > MaxLenRec { // two not-long ranges?
		goto start
	}
	shortS(ar) // we know len(ar) > MaxLenInsFC
}

// new-goroutine sort function
func glongS(ar []string, sv *syncVar) {
	longS(ar, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec
func longS(ar []string, sv *syncVar) {
start:
	aq, pv := pivotS(ar, 3)
	k := partition1S(aq, pv) // median-of-7 partitioning

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

		if len(aq) > MaxLenInsFC {
			shortS(aq)
		} else {
			insertionS(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortS(ar) // we know len(ar) > MaxLenInsFC
		return
	}

	// max goroutines? not atomic but good enough
	if sv.ngr >= MaxGor {
		longS(aq, sv) // recurse on the shorter range
		goto start
	}

	if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
		panic("sorty: longS: counter overflow")
	}
	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	go glongS(ar, sv)
	ar = aq
	goto start
}

// sortS concurrently sorts ar in ascending order.
func sortS(ar []string) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		// single-goroutine sorting
		if len(ar) > MaxLenRec {
			slongS(ar)
		} else if len(ar) > MaxLenInsFC {
			shortS(ar)
		} else if len(ar) > 1 {
			insertionS(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// median-of-9 concurrent dual partitioning with done
		k := cdualparS(ar, sv.done)
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
				panic("sorty: sortS: counter overflow")
			}
			go glongS(aq, &sv)

		} else if len(aq) > MaxLenInsFC {
			shortS(aq)
		} else {
			insertionS(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || sv.ngr >= MaxGor {
			break
		}
		// dual partition longer range
	}

	longS(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
