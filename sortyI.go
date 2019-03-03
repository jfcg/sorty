package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// int array to be sorted
var arI []int

// IsSortedI checks if ar is sorted in ascending order.
func IsSortedI(ar []int) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

func forSortI(ar []int) {
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
func ipI(pv, vl, vh int) (a, b, c int, r int) {
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
func medianI(l, h int) int {
	// lo, med, hi
	m := mean(l, h)
	vl, pv, vh := arI[l], arI[m], arI[h]

	// intermediates
	a, b := mean(l, m), mean(m, h)
	va, vb := arI[a], arI[b]

	// put lo, med, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipI(pv, vl, vh)

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipI(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipI(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipI(vh, pv, vb)
	}

	// here: vl <= va <= pv <= vb <= vh
	arI[l], arI[m], arI[h] = vl, pv, vh
	arI[a], arI[b] = va, vb
	return pv
}

var ngI, mxI uint32 // number of sorting goroutines, max limit
var doneI = make(chan bool, 1)

// SortI concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortI(ar []int, mx uint32) {
	if len(ar) < S {
		forSortI(ar)
		return
	}

	mxI = sat(mx)
	arI = ar

	ngI = 1 // count self
	gsrtI(0, len(arI)-1)
	<-doneI

	arI = nil
}

func gsrtI(lo, hi int) {
	srtI(lo, hi)

	if atomic.AddUint32(&ngI, ^uint32(0)) == 0 { // decrease goroutine counter
		doneI <- false // we are the last, all done
	}
}

// assumes hi-lo >= S-1
func srtI(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianI handles lo,hi positions

	for pv := medianI(lo, hi); l <= h; {
		swap := true
		if arI[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arI[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arI[l], arI[h] = arI[h], arI[l]
			h--
			l++
		}
	}

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= S-1 { // two big ranges?

		if ngI >= mxI { // max number of goroutines? not atomic but good enough
			srtI(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngI, 1) // increase goroutine counter
		go gsrtI(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	forSortI(arI[l : hi+1])

	if h-lo < S-1 { // two small ranges?
		forSortI(arI[lo : h+1])
		return
	}

	hi = h
	goto start
}
