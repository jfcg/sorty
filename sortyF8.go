package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// float64 array to be sorted
var arF8 []float64

// IsSortedF8 checks if ar is sorted in ascending order.
func IsSortedF8(ar []float64) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

func forSortF8(ar []float64) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

func medianF8(l, h int) float64 {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arF8[l], arF8[m], arF8[h]

	if vh < vl { // choose pivot as median of arF8[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arF8[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arF8[m] = pv
		}

		arF8[l], arF8[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arF8[m] = pv
			arF8[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arF8[m] = pv
			arF8[l] = vl
		}
	}

	return pv
}

var ngF8, mxF8 int32 // number of sorting goroutines, max limit
var doneF8 = make(chan bool, 1)

// SortF8 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting, saturated to [2, 65536].
func SortF8(ar []float64, mx int32) {
	if len(ar) < S {
		forSortF8(ar)
		return
	}

	if mx < 2 { // 2..65536 goroutines
		mxF8 = 2
	} else if mx > 65536 {
		mxF8 = 65536
	} else {
		mxF8 = mx
	}
	arF8 = ar

	ngF8 = 1 // count self
	srtF8(0, len(arF8)-1)
	<-doneF8

	arF8 = nil
}

// assumes hi-lo >= S-1
func srtF8(lo, hi int) {
	var l, h int
	var pv float64

	dec := true
	if hi < 0 { // negative hi indice means this is a recursive (slave) call
		dec = false // will not decrease counter
		hi = -hi
	}

start:
	pv = medianF8(lo, hi)
	l, h = lo+1, hi-1 // medianF8 handles lo,hi positions

	for l <= h {
		swap := true
		if arF8[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arF8[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arF8[l], arF8[h] = arF8[h], arF8[l]
			h--
			l++
		}
	}

	if hi-l < S-1 { // hi range small?
		forSortF8(arF8[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortF8(arF8[lo : h+1])

			if dec && atomic.AddInt32(&ngF8, -1) == 0 { // decrease goroutine counter
				doneF8 <- false // we are the last, all done
			}
			return // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortF8(arF8[lo : h+1])

	} else if ngF8 < mxF8 { // start a goroutine? not atomic but good enough
		atomic.AddInt32(&ngF8, 1) // increase goroutine counter
		go srtF8(lo, h)           // two big ranges, handle big lo range in another goroutine

	} else if h-lo < hi-l { // on the shorter range..
		srtF8(lo, -h) // ..start a recursive (slave) sort
	} else {
		srtF8(l, -hi)
		hi = h
		goto start
	}

	lo = l // continue with big hi range
	goto start
}
