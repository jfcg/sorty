/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "sync/atomic"

// IsSortedF8 returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedF8(ar []float64) int {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return i
		}
	}
	return 0
}

// insertion sort, assumes len(ar) >= 2
func insertionF8(ar []float64) {
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

// pivotF8 divides ar into 2n+1 equal intervals, sorts mid-points of them
// to find median-of-2n+1 pivot. ensures lo/hi ranges have at least n elements by
// moving 2n of mid-points to n positions at lo/hi ends.
// assumes n > 0, len(ar) > 4n+2. returns remaining slice,pivot for partitioning.
func pivotF8(ar []float64, n int) ([]float64, float64) {
	m := len(ar) >> 1
	s := len(ar) / (2*n + 1) // step > 1
	l, h := m-n*s, m+n*s

	for q, k := h, m-2*s; k >= l; { // insertion sort ar[m+i*s], i=-n..n
		if ar[q] < ar[k] {
			ar[k], ar[q] = ar[q], ar[k]
		}
		q -= s
		k -= s
	}
	for q := l; ; {
		k := q
		q += s
		v := ar[q]
		if v < ar[k] {
			for {
				ar[k+s] = ar[k]
				k -= s
				if k < l || v >= ar[k] {
					break
				}
			}
			ar[k+s] = v
		}
		if q >= h {
			break
		}
	}

	lo, hi := 0, len(ar)-1

	// move lo/hi mid-points to lo/hi ends
	for {
		ar[l], ar[lo] = ar[lo], ar[l]
		ar[h], ar[hi] = ar[hi], ar[h]
		l += s
		h -= s
		lo++
		if h <= m {
			break
		}
		hi--
	}

	return ar[lo:hi:hi], ar[m] // lo <= m-s+1, m+s-1 < hi
}

// partition ar into <= and >= pivot, assumes len(ar) >= 2
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func partition1F8(ar []float64, pv float64) int {
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
func partition2F8(ar []float64, a, b int, pv float64) (int, int) {
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

// concurrent dual partitioning of ar
// returns short & long sub-ranges
func cdualparF8(ar []float64, ch chan int) (s, l []float64) {

	aq, pv := pivotF8(ar, 4) // median-of-9
	k := len(aq) >> 1
	a, b := k>>1, mid(k, len(aq))

	go func(ap []float64) {
		ch <- partition1F8(ap, pv) // mid half range
	}(aq[a:b:b])

	t := a
	a, b = partition2F8(aq, a, b, pv) // left/right quarter ranges
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

	k += 4 // convert k indice to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}
	return aq, ar
}

// partitions ar with uniform median-of-2n+1 pivot and
// returns short & long sub-ranges
func partF8(ar []float64, n int) (s, l []float64) {

	aq, pv := pivotF8(ar, n)
	k := partition1F8(aq, pv)

	k += n // convert k indice from aq to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}
	return aq, ar
}

// short range sort function, assumes Mli < len(ar) <= Mlr
func shortF8(ar []float64) {
start:
	aq, ar := partF8(ar, 2) // median-of-5 partitioning

	if len(aq) > Mli {
		shortF8(aq) // recurse on the shorter range
		goto start
	}
	insertionF8(aq) // at least one insertion range

	if len(ar) > Mli {
		goto start
	}
	insertionF8(ar) // two insertion ranges
	return
}

// SortF8 concurrently sorts ar in ascending order.
func SortF8(ar []float64) {
	var (
		ngr  = uint32(1)     // number of sorting goroutines including this
		done chan int        // end signal
		long func([]float64) // long range sort function
	)

	glong := func(ar []float64) { // new-goroutine sort function
		long(ar)
		if atomic.AddUint32(&ngr, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- 0 // we are the last, all done
		}
	}

	long = func(ar []float64) { // assumes len(ar) > Mlr
	start:
		aq, ar := partF8(ar, 3) // median-of-7 partitioning

		// branches below are optimal for fewer total jumps
		if len(aq) <= Mlr { // at least one not-long range?

			if len(aq) > Mli {
				shortF8(aq)
			} else {
				insertionF8(aq)
			}

			if len(ar) > Mlr { // two not-long ranges?
				goto start
			}
			shortF8(ar) // we know len(ar) > Mli
			return
		}

		// max goroutines? not atomic but good enough
		if ngr >= Mxg {
			long(aq) // recurse on the shorter range
			goto start
		}

		if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
			panic("SortF8: long: counter overflow")
		}
		// new-goroutine sort on the longer range only when
		// both ranges are big and max goroutines is not exceeded
		go glong(ar)
		ar = aq
		goto start
	}

	if len(ar) < 2*(Mlr+1) {
		if len(ar) > Mlr {
			long(ar) // will not create goroutines or use ngr/done

		} else if len(ar) > Mli {
			shortF8(ar)
		} else if len(ar) > 1 {
			insertionF8(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	done = make(chan int, 1) // maybe this goroutine will be the last
	for {
		// median-of-9 concurrent dual partitioning with done
		var aq []float64
		aq, ar = cdualparF8(ar, done)

		// handle shorter range
		if len(aq) > Mlr {
			if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
				panic("SortF8: dual: counter overflow")
			}
			go glong(aq)

		} else if len(aq) > Mli {
			shortF8(aq)
		} else {
			insertionF8(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(Mlr+1) || ngr >= Mxg {
			break
		}
		// dual partition longer range
	}

	glong(ar) // we know len(ar) > Mlr
	<-done
}
