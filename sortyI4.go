package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// int32 array to be sorted
var arI4 []int32

// IsSortedI4 checks if ar is sorted in ascending order.
func IsSortedI4(ar []int32) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

func forSortI4(ar []int32) {
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
func ipI4(pv, vl, vh int32) (a, b, c int32, r int) {
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
func medianI4(l, h int) int32 {
	// lo, med, hi
	m := mean(l, h)
	vl, pv, vh := arI4[l], arI4[m], arI4[h]

	// intermediates
	a, b := mean(l, m), mean(m, h)
	va, vb := arI4[a], arI4[b]

	// put lo, med, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipI4(pv, vl, vh)

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipI4(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipI4(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipI4(vh, pv, vb)
	}

	// here: vl <= va <= pv <= vb <= vh
	arI4[l], arI4[m], arI4[h] = vl, pv, vh
	arI4[a], arI4[b] = va, vb
	return pv
}

var ngI4, mxI4 uint32 // number of sorting goroutines, max limit
var doneI4 = make(chan bool, 1)

// SortI4 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortI4(ar []int32, mx uint32) {
	if len(ar) < S {
		forSortI4(ar)
		return
	}

	mxI4 = sat(mx)
	arI4 = ar

	ngI4 = 1 // count self
	gsrtI4(0, len(arI4)-1)
	<-doneI4

	arI4 = nil
}

func gsrtI4(lo, hi int) {
	srtI4(lo, hi)

	if atomic.AddUint32(&ngI4, ^uint32(0)) == 0 { // decrease goroutine counter
		doneI4 <- false // we are the last, all done
	}
}

// assumes hi-lo >= S-1
func srtI4(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianI4 handles lo,hi positions

	for pv := medianI4(lo, hi); l <= h; {
		swap := true
		if arI4[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arI4[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arI4[l], arI4[h] = arI4[h], arI4[l]
			h--
			l++
		}
	}

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= S-1 { // two big ranges?

		if ngI4 >= mxI4 { // max number of goroutines? not atomic but good enough
			srtI4(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngI4, 1) // increase goroutine counter
		go gsrtI4(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	forSortI4(arI4[l : hi+1])

	if h-lo < S-1 { // two small ranges?
		forSortI4(arI4[lo : h+1])
		return
	}

	hi = h
	goto start
}
