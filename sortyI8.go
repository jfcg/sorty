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

func medianI8(l, h int) int64 {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arI8[l], arI8[m], arI8[h]

	if vh < vl { // choose pivot as median of arI8[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arI8[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arI8[m] = pv
		}

		arI8[l], arI8[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arI8[m] = pv
			arI8[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arI8[m] = pv
			arI8[l] = vl
		}
	}

	return pv
}

var ngI8, mxI8 int32 // number of sorting goroutines, max limit
var doneI8 = make(chan bool, 1)

// SortI8 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting, saturated to [2, 65536].
func SortI8(ar []int64, mx int32) {
	if len(ar) < S {
		forSortI8(ar)
		return
	}

	if mx < 2 { // 2..65536 goroutines
		mxI8 = 2
	} else if mx > 65536 {
		mxI8 = 65536
	} else {
		mxI8 = mx
	}
	arI8 = ar

	ngI8 = 1 // count self
	srtI8(0, len(arI8)-1)
	<-doneI8

	arI8 = nil
}

// assumes hi-lo >= S-1
func srtI8(lo, hi int) {
	var l, h int
	var pv int64

	dec := true
	if hi < 0 { // negative hi indice means this is a recursive (slave) call
		dec = false // will not decrease counter
		hi = -hi
	}

start:
	pv = medianI8(lo, hi)
	l, h = lo+1, hi-1 // medianI8 handles lo,hi positions

	for l <= h {
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

	if hi-l < S-1 { // hi range small?
		forSortI8(arI8[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortI8(arI8[lo : h+1])

			if dec && atomic.AddInt32(&ngI8, -1) == 0 { // decrease goroutine counter
				doneI8 <- false // we are the last, all done
			}
			return // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortI8(arI8[lo : h+1])

	} else if ngI8 < mxI8 { // start a goroutine? not atomic but good enough
		atomic.AddInt32(&ngI8, 1) // increase goroutine counter
		go srtI8(lo, h)           // two big ranges, handle big lo range in another goroutine

	} else if h-lo < hi-l { // on the shorter range..
		srtI8(lo, -h) // ..start a recursive (slave) sort
	} else {
		srtI8(l, -hi)
		hi = h
		goto start
	}

	lo = l // continue with big hi range
	goto start
}
