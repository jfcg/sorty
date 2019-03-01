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

func medianP(l, h int) uintptr {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arP[l], arP[m], arP[h]

	if vh < vl { // choose pivot as median of arP[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arP[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arP[m] = pv
		}

		arP[l], arP[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arP[m] = pv
			arP[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arP[m] = pv
			arP[l] = vl
		}
	}

	return pv
}

var ngP, mxP int32 // number of sorting goroutines, max limit
var doneP = make(chan bool, 1)

// SortP concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting, saturated to [2, 65536].
func SortP(ar []uintptr, mx int32) {
	if len(ar) < S {
		forSortP(ar)
		return
	}

	if mx < 2 { // 2..65536 goroutines
		mxP = 2
	} else if mx > 65536 {
		mxP = 65536
	} else {
		mxP = mx
	}
	arP = ar

	ngP = 1 // count self
	srtP(0, len(arP)-1)
	<-doneP

	arP = nil
}

// assumes hi-lo >= S-1
func srtP(lo, hi int) {
	var l, h int
	var pv uintptr

	dec := true
	if hi < 0 { // negative hi indice means this is a recursive (slave) call
		dec = false // will not decrease counter
		hi = -hi
	}

start:
	pv = medianP(lo, hi)
	l, h = lo+1, hi-1 // medianP handles lo,hi positions

	for l <= h {
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

	if hi-l < S-1 { // hi range small?
		forSortP(arP[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortP(arP[lo : h+1])

			if dec && atomic.AddInt32(&ngP, -1) == 0 { // decrease goroutine counter
				doneP <- false // we are the last, all done
			}
			return // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortP(arP[lo : h+1])

	} else if ngP < mxP { // start a goroutine? not atomic but good enough
		atomic.AddInt32(&ngP, 1) // increase goroutine counter
		go srtP(lo, h)           // two big ranges, handle big lo range in another goroutine

	} else if h-lo < hi-l { // on the shorter range..
		srtP(lo, -h) // ..start a recursive (slave) sort
	} else {
		srtP(l, -hi)
		hi = h
		goto start
	}

	lo = l // continue with big hi range
	goto start
}
