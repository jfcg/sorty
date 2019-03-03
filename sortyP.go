package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// uintptr array to be sorted
var arP []uintptr

// IsSortedP checks if ar is sorted in ascending order.
func IsSortedP(ar []uintptr) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

func forSortP(ar []uintptr) {
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
func ipP(pv, vl, vh uintptr) (a, b, c uintptr, r int) {
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
func medianP(l, h int) uintptr {
	// lo, med, hi
	m := mean(l, h)
	vl, pv, vh := arP[l], arP[m], arP[h]

	// intermediates
	a, b := mean(l, m), mean(m, h)
	va, vb := arP[a], arP[b]

	// put lo, med, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipP(pv, vl, vh)

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipP(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipP(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipP(vh, pv, vb)
	}

	// here: vl <= va <= pv <= vb <= vh
	arP[l], arP[m], arP[h] = vl, pv, vh
	arP[a], arP[b] = va, vb
	return pv
}

var ngP, mxP uint32 // number of sorting goroutines, max limit
var doneP = make(chan bool, 1)

// SortP concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortP(ar []uintptr, mx uint32) {
	if len(ar) < S {
		forSortP(ar)
		return
	}

	mxP = sat(mx)
	arP = ar

	ngP = 1 // count self
	gsrtP(0, len(arP)-1)
	<-doneP

	arP = nil
}

func gsrtP(lo, hi int) {
	srtP(lo, hi)

	if atomic.AddUint32(&ngP, ^uint32(0)) == 0 { // decrease goroutine counter
		doneP <- false // we are the last, all done
	}
}

// assumes hi-lo >= S-1
func srtP(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianP handles lo,hi positions

	for pv := medianP(lo, hi); l <= h; {
		swap := true
		if arP[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arP[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arP[l], arP[h] = arP[h], arP[l]
			h--
			l++
		}
	}

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= S-1 { // two big ranges?

		if ngP >= mxP { // max number of goroutines? not atomic but good enough
			srtP(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngP, 1) // increase goroutine counter
		go gsrtP(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	forSortP(arP[l : hi+1])

	if h-lo < S-1 { // two small ranges?
		forSortP(arP[lo : h+1])
		return
	}

	hi = h
	goto start
}
