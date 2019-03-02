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

func medianS(l, h int) string {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arS[l], arS[m], arS[h]

	if vh < vl { // choose pivot as median of arS[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arS[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arS[m] = pv
		}

		arS[l], arS[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arS[m] = pv
			arS[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arS[m] = pv
			arS[l] = vl
		}
	}

	return pv
}

var ngS, mxS uint32 // number of sorting goroutines, max limit
var doneS = make(chan bool, 1)

// SortS concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting, saturated to [2, 65536].
func SortS(ar []string, mx uint32) {
	if len(ar) < S {
		forSortS(ar)
		return
	}

	if mx < 2 { // 2..65536 goroutines
		mxS = 2
	} else if mx > 65536 {
		mxS = 65536
	} else {
		mxS = mx
	}
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
