/*	Copyright (c) 2021, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"sync/atomic"

	"github.com/jfcg/sixb"
)

// IsSortedB returns 0 if ar is sorted in ascending lexicographical order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedB(ar [][]byte) int {
	for i := len(ar) - 1; i > 0; i-- {
		if sixb.BtoS(ar[i]) < sixb.BtoS(ar[i-1]) {
			return i
		}
	}
	return 0
}

// insertion sort, assumes len(ar) >= 2
func insertionB(ar [][]byte) {
	hi := len(ar) - 1
	for l, h := (hi-3)>>1, hi; l >= 0; {
		if sixb.BtoS(ar[h]) < sixb.BtoS(ar[l]) {
			ar[l], ar[h] = ar[h], ar[l]
		}
		l--
		h--
	}
	for h := 0; ; {
		l := h
		h++
		x := ar[h]
		v := sixb.BtoS(x)
		if v < sixb.BtoS(ar[l]) {
			for {
				ar[l+1] = ar[l]
				l--
				if l < 0 || v >= sixb.BtoS(ar[l]) {
					break
				}
			}
			ar[l+1] = x
		}
		if h >= hi {
			break
		}
	}
}

// pivotB divides ar into 2n+1 equal intervals, sorts mid-points of them
// to find median-of-2n+1 pivot. ensures lo/hi ranges have at least n elements by
// moving 2n of mid-points to n positions at lo/hi ends.
// assumes n > 0, len(ar) > 4n+2. returns remaining slice,pivot for partitioning.
func pivotB(ar [][]byte, n int) ([][]byte, string) {
	m := len(ar) >> 1
	s := len(ar) / (2*n + 1) // step > 1
	l, h := m-n*s, m+n*s

	for q, k := h, m-2*s; k >= l; { // insertion sort ar[m+i*s], i=-n..n
		if sixb.BtoS(ar[q]) < sixb.BtoS(ar[k]) {
			ar[k], ar[q] = ar[q], ar[k]
		}
		q -= s
		k -= s
	}
	for q := l; ; {
		k := q
		q += s
		x := ar[q]
		v := sixb.BtoS(x)
		if v < sixb.BtoS(ar[k]) {
			for {
				ar[k+s] = ar[k]
				k -= s
				if k < l || v >= sixb.BtoS(ar[k]) {
					break
				}
			}
			ar[k+s] = x
		}
		if q >= h {
			break
		}
	}

	lo, hi := 0, len(ar)

	// move lo/hi mid-points to lo/hi ends
	for {
		hi--
		ar[l], ar[lo] = ar[lo], ar[l]
		ar[h], ar[hi] = ar[hi], ar[h]
		l += s
		h -= s
		lo++
		if h <= m {
			break
		}
	}

	return ar[lo:hi:hi], sixb.BtoS(ar[m]) // lo <= m-s+1, m+s-1 < hi
}

// partition ar into <= and >= pivot, assumes len(ar) >= 2
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func partition1B(ar [][]byte, pv string) int {
	l, h := 0, len(ar)-1
	for {
		if sixb.BtoS(ar[h]) < pv { // avoid unnecessary comparisons
			for {
				if pv < sixb.BtoS(ar[l]) {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if pv < sixb.BtoS(ar[l]) { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l
				}
				if sixb.BtoS(ar[h]) < pv {
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
	if l == h && sixb.BtoS(ar[h]) < pv { // classify mid element
		l++
	}
	return l
}

// rearrange ar[:a] and ar[b:] into <= and >= pivot, assumes 0 < a < b < len(ar)
// gap (a,b) expands until one of the intervals is fully consumed
func partition2B(ar [][]byte, a, b int, pv string) (int, int) {
	a--
	for {
		if sixb.BtoS(ar[b]) < pv { // avoid unnecessary comparisons
			for {
				if pv < sixb.BtoS(ar[a]) {
					ar[a], ar[b] = ar[b], ar[a]
					break
				}
				a--
				if a < 0 {
					return a, b
				}
			}
		} else if pv < sixb.BtoS(ar[a]) { // extend ranges in balance
			for {
				b++
				if b >= len(ar) {
					return a, b
				}
				if sixb.BtoS(ar[b]) < pv {
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
func gpart1B(ar [][]byte, pv string, ch chan int) {
	ch <- partition1B(ar, pv)
}

// concurrent dual partitioning of ar
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func cdualparB(ar [][]byte, ch chan int) int {

	aq, pv := pivotB(ar, 4) // median-of-9
	k := len(aq) >> 1
	a, b := k>>1, mid(k, len(aq))

	go gpart1B(aq[a:b:b], pv, ch) // mid half range

	t := a
	a, b = partition2B(aq, a, b, pv) // left/right quarter ranges
	k = <-ch
	k += t // convert k indice to aq

	// only one gap is possible
	for ; 0 <= a; a-- { // gap left in low range?
		if pv < sixb.BtoS(aq[a]) {
			k--
			aq[a], aq[k] = aq[k], aq[a]
		}
	}
	for ; b < len(aq); b++ { // gap left in high range?
		if sixb.BtoS(aq[b]) < pv {
			aq[b], aq[k] = aq[k], aq[b]
			k++
		}
	}
	return k + 4 // convert k indice to ar
}

// short range sort function, assumes Hmli < len(ar) <= Mlr
func shortB(ar [][]byte) {
start:
	aq, pv := pivotB(ar, 2)
	k := partition1B(aq, pv) // median-of-5 partitioning

	k += 2 // convert k indice from aq to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > Hmli {
		shortB(aq) // recurse on the shorter range
		goto start
	}
	insertionB(aq) // at least one insertion range

	if len(ar) > Hmli {
		goto start
	}
	insertionB(ar) // two insertion ranges
}

// long range sort function (single goroutine), assumes len(ar) > Mlr
func slongB(ar [][]byte) {
start:
	aq, pv := pivotB(ar, 3)
	k := partition1B(aq, pv) // median-of-7 partitioning

	k += 3 // convert k indice from aq to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > Mlr { // at least one not-long range?
		slongB(aq) // recurse on the shorter range
		goto start
	}

	if len(aq) > Hmli {
		shortB(aq)
	} else {
		insertionB(aq)
	}

	if len(ar) > Mlr { // two not-long ranges?
		goto start
	}
	shortB(ar) // we know len(ar) > Hmli
}

// new-goroutine sort function
func glongB(ar [][]byte, sv *syncVar) {
	longB(ar, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > Mlr
func longB(ar [][]byte, sv *syncVar) {
start:
	aq, pv := pivotB(ar, 3)
	k := partition1B(aq, pv) // median-of-7 partitioning

	k += 3 // convert k indice from aq to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	// branches below are optimal for fewer total jumps
	if len(aq) <= Mlr { // at least one not-long range?

		if len(aq) > Hmli {
			shortB(aq)
		} else {
			insertionB(aq)
		}

		if len(ar) > Mlr { // two not-long ranges?
			goto start
		}
		shortB(ar) // we know len(ar) > Hmli
		return
	}

	// max goroutines? not atomic but good enough
	if sv.ngr >= Mxg {
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

// SortB concurrently sorts ar in ascending lexicographical order.
func SortB(ar [][]byte) {

	if len(ar) < 2*(Mlr+1) || Mxg <= 1 {

		// single-goroutine sorting
		if len(ar) > Mlr {
			slongB(ar)
		} else if len(ar) > Hmli {
			shortB(ar)
		} else if len(ar) > 1 {
			insertionB(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// median-of-9 concurrent dual partitioning with done
		k := cdualparB(ar, sv.done)
		var aq [][]byte

		if k < len(ar)-k {
			aq = ar[:k:k]
			ar = ar[k:] // ar is the longer range
		} else {
			aq = ar[k:]
			ar = ar[:k:k]
		}

		// handle shorter range
		if len(aq) > Mlr {
			if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
				panic("sorty: SortB: counter overflow")
			}
			go glongB(aq, &sv)

		} else if len(aq) > Hmli {
			shortB(aq)
		} else {
			insertionB(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(Mlr+1) || sv.ngr >= Mxg {
			break
		}
		// dual partition longer range
	}

	longB(ar, &sv) // we know len(ar) > Mlr

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
