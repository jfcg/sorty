package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// uint32 array to be sorted
var arU4 []uint32

// IsSortedU4 checks if ar is sorted in ascending order.
func IsSortedU4(ar []uint32) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

func forSortU4(ar []uint32) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

func medianU4(l, h int) uint32 {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arU4[l], arU4[m], arU4[h]

	if vh < vl { // choose pivot as median of arU4[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arU4[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arU4[m] = pv
		}

		arU4[l], arU4[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arU4[m] = pv
			arU4[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arU4[m] = pv
			arU4[l] = vl
		}
	}

	return pv
}

var ngU4, mxU4 int32 // number of sorting goroutines, max limit
var doneU4 = make(chan bool, 1)

// SortU4 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting, saturated to [2, 65536].
func SortU4(ar []uint32, mx int32) {
	if len(ar) < S {
		forSortU4(ar)
		return
	}

	if mx < 2 { // 2..65536 goroutines
		mxU4 = 2
	} else if mx > 65536 {
		mxU4 = 65536
	} else {
		mxU4 = mx
	}
	arU4 = ar

	ngU4 = 1 // count self
	srtU4(0, len(arU4)-1)
	<-doneU4

	arU4 = nil
}

// assumes hi-lo >= S-1
func srtU4(lo, hi int) {
	var l, h int
	var pv uint32

	dec := true
	if hi < 0 { // negative hi indice means this is a recursive (slave) call
		dec = false // will not decrease counter
		hi = -hi
	}

start:
	pv = medianU4(lo, hi)
	l, h = lo+1, hi-1 // medianU4 handles lo,hi positions

	for l <= h {
		swap := true
		if arU4[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arU4[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arU4[l], arU4[h] = arU4[h], arU4[l]
			h--
			l++
		}
	}

	if hi-l < S-1 { // hi range small?
		forSortU4(arU4[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortU4(arU4[lo : h+1])

			if dec && atomic.AddInt32(&ngU4, -1) == 0 { // decrease goroutine counter
				doneU4 <- false // we are the last, all done
			}
			return // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortU4(arU4[lo : h+1])

	} else if ngU4 < mxU4 { // start a goroutine? not atomic but good enough
		atomic.AddInt32(&ngU4, 1) // increase goroutine counter
		go srtU4(lo, h)           // two big ranges, handle big lo range in another goroutine

	} else if h-lo < hi-l { // on the shorter range..
		srtU4(lo, -h) // ..start a recursive (slave) sort
	} else {
		srtU4(l, -hi)
		hi = h
		goto start
	}

	lo = l // continue with big hi range
	goto start
}
