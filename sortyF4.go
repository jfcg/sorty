/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "sync/atomic"

// IsSortedF4 returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedF4(ar []float32) int {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return i
		}
	}
	return 0
}

// insertion sort, assumes 0 < hi < len(ar)
func insertionF4(ar []float32, hi int) {

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

// pivotF4 divides [lo..hi] range into 2n+1 equal intervals, sorts mid-points of them
// to find median-of-2n+1 pivot. ensures lo/hi ranges have at least n elements by
// moving 2n of mid-points to n positions at lo/hi ends.
// assumes n > 0, lo+4n+1 < hi. returns start,pivot,end for partitioning.
func pivotF4(ar []float32, lo, hi, n int) (int, float32, int) {
	m := mid(lo, hi)
	s := int(uint(hi-lo+1) / uint(2*n+1)) // step > 1
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

	// move hi mid-points to hi end
	for {
		ar[h], ar[hi] = ar[hi], ar[h]
		h -= s
		hi--
		if h <= m {
			break
		}
	}

	// move lo mid-points to lo end
	for {
		ar[l], ar[lo] = ar[lo], ar[l]
		l += s
		lo++
		if l >= m {
			break
		}
	}
	return lo, ar[m], hi // lo <= m-s+1, m+s-1 <= hi
}

// partition ar[l..h] into <= and >= pivot, assumes l < h
// returns m with ar[:m] <= pivot, ar[m:] >= pivot
func partition1F4(ar []float32, l int, pv float32, h int) int {
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

// rearrange ar[l..a] and ar[b..h] into <= and >= pivot, assumes l <= a < b <= h
// gap (a..b) expands until one of the intervals is fully consumed
func partition2F4(ar []float32, l, a int, pv float32, b, h int) (int, int) {
	for {
		if ar[b] < pv { // avoid unnecessary comparisons
			for {
				if pv < ar[a] {
					ar[a], ar[b] = ar[b], ar[a]
					break
				}
				a--
				if a < l {
					return a, b
				}
			}
		} else if pv < ar[a] { // extend ranges in balance
			for {
				b++
				if b > h {
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
		if a < l || b > h {
			return a, b
		}
	}
}

// concurrent dual partitioning
// returns m with ar[:m] <= pivot, ar[m:] >= pivot
func cdualparF4(par chan int, ar []float32, lo, hi int) int {

	lo, pv, hi := pivotF4(ar, lo, hi, 4) // median-of-9

	m := mid(lo, hi)
	a, b := mid(lo, m), mid(m, hi)

	go func(l, h int) {
		par <- partition1F4(ar, l, pv, h) // mid half range
	}(a, b)

	a, b = partition2F4(ar, lo, a-1, pv, b+1, hi) // left/right quarter ranges
	m = <-par

	// only one gap is possible
	for ; lo <= a; a-- { // gap left in low range?
		if pv < ar[a] {
			m--
			ar[a], ar[m] = ar[m], ar[a]
		}
	}
	for ; b <= hi; b++ { // gap left in high range?
		if ar[b] < pv {
			ar[b], ar[m] = ar[m], ar[b]
			m++
		}
	}
	return m
}

// short range sort function, assumes Mli <= no < Mlr
func shortF4(ar []float32, lo, no int) {
start:
	n := lo + no
	l, pv, h := pivotF4(ar, lo, n, 2) // median-of-5
	l = partition1F4(ar, l, pv, h)
	n -= l
	no -= n + 1

	if no < n {
		n, no = no, n // [lo,lo+no] is the longer range
		l, lo = lo, l
	}

	if n >= Mli {
		shortF4(ar, l, n) // recurse on the shorter range
		goto start
	}
	insertionF4(ar[l:], n) // at least one insertion range

	if no >= Mli {
		goto start
	}
	insertionF4(ar[lo:], no) // two insertion ranges
	return
}

// SortF4 concurrently sorts ar in ascending order.
func SortF4(ar []float32) {
	var (
		ngr  = uint32(1)    // number of sorting goroutines including this
		done chan int       // end signal
		long func(int, int) // long range sort function
	)

	glong := func(lo, no int) { // new-goroutine sort function
		long(lo, no)
		if atomic.AddUint32(&ngr, ^uint32(0)) == 0 { // decrease goroutine counter
			done <- 0 // we are the last, all done
		}
	}

	long = func(lo, no int) { // assumes no >= Mlr
	start:
		n := lo + no
		l, pv, h := pivotF4(ar, lo, n, 3) // median-of-7
		l = partition1F4(ar, l, pv, h)
		n -= l
		no -= n + 1

		if no < n {
			n, no = no, n // [lo,lo+no] is the longer range
			l, lo = lo, l
		}

		// branches below are optimal for fewer total jumps
		if n < Mlr { // at least one not-long range?
			if n >= Mli {
				shortF4(ar, l, n)
			} else {
				insertionF4(ar[l:], n)
			}

			if no >= Mlr { // two not-long ranges?
				goto start
			}
			shortF4(ar, lo, no) // we know no >= Mli
			return
		}

		// max goroutines? not atomic but good enough
		if ngr >= Mxg {
			long(l, n) // recurse on the shorter range
			goto start
		}

		if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
			panic("SortF4: long: counter overflow")
		}
		// new-goroutine sort on the longer range only when
		// both ranges are big and max goroutines is not exceeded
		go glong(lo, no)
		lo, no = l, n
		goto start
	}

	n := len(ar) - 1 // high indice
	if n <= 2*Mlr {
		if n >= Mlr {
			long(0, n) // will not create goroutines or use ngr/done
		} else if n >= Mli {
			shortF4(ar, 0, n)
		} else if n > 0 {
			insertionF4(ar, n)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	done = make(chan int, 1) // maybe this goroutine will be the last
	lo, no := 0, n
	for {
		// concurrent dual partitioning with done
		l := cdualparF4(done, ar, lo, n)
		n -= l
		no -= n + 1

		if no < n {
			n, no = no, n // [lo,lo+no] is the longer range
			l, lo = lo, l
		}

		// handle shorter range
		if n >= Mlr {
			if atomic.AddUint32(&ngr, 1) == 0 { // increase goroutine counter
				panic("SortF4: dual: counter overflow")
			}
			go glong(l, n)

		} else if n >= Mli {
			shortF4(ar, l, n)
		} else {
			insertionF4(ar[l:], n)
		}

		// longer range big enough? max goroutines?
		if no <= 2*Mlr || ngr >= Mxg {
			break
		}
		n = lo + no // dual partition longer range
	}

	glong(lo, no) // we know no >= Mlr
	<-done
}
