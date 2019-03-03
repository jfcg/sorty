package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// string array to be sorted
var arS []string

// IsSortedS checks if ar is sorted in ascending order.
func IsSortedS(ar []string) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

func forSortS(ar []string) {
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
func ipS(pv, vl, vh string) (a, b, c string, r int) {
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
func medianS(l, h int) string {
	// lo, med, hi
	m := mean(l, h)
	vl, pv, vh := arS[l], arS[m], arS[h]

	// intermediates
	a, b := mean(l, m), mean(m, h)
	va, vb := arS[a], arS[b]

	// put lo, med, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipS(pv, vl, vh)

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipS(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipS(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipS(vh, pv, vb)
	}

	// here: vl <= va <= pv <= vb <= vh
	arS[l], arS[m], arS[h] = vl, pv, vh
	arS[a], arS[b] = va, vb
	return pv
}

var ngS, mxS uint32 // number of sorting goroutines, max limit
var doneS = make(chan bool, 1)

// SortS concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortS(ar []string, mx uint32) {
	if len(ar) < S {
		forSortS(ar)
		return
	}

	mxS = sat(mx)
	arS = ar

	ngS = 1 // count self
	gsrtS(0, len(arS)-1)
	<-doneS

	arS = nil
}

func gsrtS(lo, hi int) {
	srtS(lo, hi)

	if atomic.AddUint32(&ngS, ^uint32(0)) == 0 { // decrease goroutine counter
		doneS <- false // we are the last, all done
	}
}

// assumes hi-lo >= S-1
func srtS(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianS handles lo,hi positions

	for pv := medianS(lo, hi); l <= h; {
		swap := true
		if arS[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arS[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arS[l], arS[h] = arS[h], arS[l]
			h--
			l++
		}
	}

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= S-1 { // two big ranges?

		if ngS >= mxS { // max number of goroutines? not atomic but good enough
			srtS(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngS, 1) // increase goroutine counter
		go gsrtS(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	forSortS(arS[l : hi+1])

	if h-lo < S-1 { // two small ranges?
		forSortS(arS[lo : h+1])
		return
	}

	hi = h
	goto start
}
