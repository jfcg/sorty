package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// int64 array to be sorted
var arI8 []int64

// IsSortedI8 checks if ar is sorted in ascending order.
func IsSortedI8(ar []int64) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

func forSortI8(ar []int64) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

// given vl <= vh, inserts pv in the middle
// returns vl <= pv <= vh
func ipI8(pv, vl, vh int64) (a, b, c int64, r int) {
	if pv > vh {
		vh, pv = pv, vh
		r = 1
	} else if pv < vl {
		vl, pv = pv, vl
		r = -1
	}
	return vl, pv, vh, r
}

// return pivot as median of five scattered values
func medianI8(l, h int) int64 {
	// lo, med, hi
	m := mean(l, h)
	vl, pv, vh := arI8[l], arI8[m], arI8[h]

	// intermediates
	a, b := mean(l, m), mean(m, h)
	va, vb := arI8[a], arI8[b]

	// put lo, med, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipI8(pv, vl, vh)

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipI8(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipI8(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipI8(vh, pv, vb)
	}

	// here: vl <= va <= pv <= vb <= vh
	arI8[l], arI8[m], arI8[h] = vl, pv, vh
	arI8[a], arI8[b] = va, vb
	return pv
}

var ngI8, mxI8 uint32 // number of sorting goroutines, max limit
var doneI8 = make(chan bool, 1)

// SortI8 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortI8(ar []int64, mx uint32) {
	if len(ar) < S {
		forSortI8(ar)
		return
	}

	mxI8 = sat(mx)
	arI8 = ar

	ngI8 = 1 // count self
	gsrtI8(0, len(arI8)-1)
	<-doneI8

	arI8 = nil
}

func gsrtI8(lo, hi int) {
	srtI8(lo, hi)

	if atomic.AddUint32(&ngI8, ^uint32(0)) == 0 { // decrease goroutine counter
		doneI8 <- false // we are the last, all done
	}
}

// assumes hi-lo >= S-1
func srtI8(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianI8 handles lo,hi positions

	for pv := medianI8(lo, hi); l <= h; {
		swap := true
		if arI8[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arI8[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arI8[l], arI8[h] = arI8[h], arI8[l]
			h--
			l++
		}
	}

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= S-1 { // two big ranges?

		if ngI8 >= mxI8 { // max number of goroutines? not atomic but good enough
			srtI8(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngI8, 1) // increase goroutine counter
		go gsrtI8(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	forSortI8(arI8[l : hi+1])

	if h-lo < S-1 { // two small ranges?
		forSortI8(arI8[lo : h+1])
		return
	}

	hi = h
	goto start
}
