package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// float32 array to be sorted
var arF4 []float32

// IsSortedF4 checks if ar is sorted in ascending order.
func IsSortedF4(ar []float32) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

func forSortF4(ar []float32) {
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
func ipF4(pv, vl, vh float32) (a, b, c float32, r int) {
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
func medianF4(l, h int) float32 {
	// lo, med, hi
	m := mean(l, h)
	vl, pv, vh := arF4[l], arF4[m], arF4[h]

	// intermediates
	a, b := mean(l, m), mean(m, h)
	va, vb := arF4[a], arF4[b]

	// put lo, med, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipF4(pv, vl, vh)

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipF4(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipF4(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipF4(vh, pv, vb)
	}

	// here: vl <= va <= pv <= vb <= vh
	arF4[l], arF4[m], arF4[h] = vl, pv, vh
	arF4[a], arF4[b] = va, vb
	return pv
}

var ngF4, mxF4 uint32 // number of sorting goroutines, max limit
var doneF4 = make(chan bool, 1)

// SortF4 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortF4(ar []float32, mx uint32) {
	if len(ar) < S {
		forSortF4(ar)
		return
	}

	mxF4 = sat(mx)
	arF4 = ar

	ngF4 = 1 // count self
	gsrtF4(0, len(arF4)-1)
	<-doneF4

	arF4 = nil
}

func gsrtF4(lo, hi int) {
	srtF4(lo, hi)

	if atomic.AddUint32(&ngF4, ^uint32(0)) == 0 { // decrease goroutine counter
		doneF4 <- false // we are the last, all done
	}
}

// assumes hi-lo >= S-1
func srtF4(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianF4 handles lo,hi positions

	for pv := medianF4(lo, hi); l <= h; {
		swap := true
		if arF4[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arF4[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arF4[l], arF4[h] = arF4[h], arF4[l]
			h--
			l++
		}
	}

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= S-1 { // two big ranges?

		if ngF4 >= mxF4 { // max number of goroutines? not atomic but good enough
			srtF4(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngF4, 1) // increase goroutine counter
		go gsrtF4(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	forSortF4(arF4[l : hi+1])

	if h-lo < S-1 { // two small ranges?
		forSortF4(arF4[lo : h+1])
		return
	}

	hi = h
	goto start
}
