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

func medianI4(l, h int) int32 {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arI4[l], arI4[m], arI4[h]

	if vh < vl { // choose pivot as median of arI4[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arI4[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arI4[m] = pv
		}

		arI4[l], arI4[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arI4[m] = pv
			arI4[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arI4[m] = pv
			arI4[l] = vl
		}
	}

	return pv
}

var ngI4, mxI4 uint32 // number of sorting goroutines, max limit
var doneI4 = make(chan bool, 1)

// SortI4 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting, saturated to [2, 65536].
func SortI4(ar []int32, mx uint32) {
	if len(ar) < S {
		forSortI4(ar)
		return
	}

	if mx < 2 { // 2..65536 goroutines
		mxI4 = 2
	} else if mx > 65536 {
		mxI4 = 65536
	} else {
		mxI4 = mx
	}
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
