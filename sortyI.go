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

func medianI(l, h int) int {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arI[l], arI[m], arI[h]

	if vh < vl { // choose pivot as median of arI[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arI[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arI[m] = pv
		}

		arI[l], arI[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arI[m] = pv
			arI[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arI[m] = pv
			arI[l] = vl
		}
	}

	return pv
}

var ngI, mxI uint32 // number of sorting goroutines, max limit
var doneI = make(chan bool, 1)

// SortI concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting, saturated to [2, 65536].
func SortI(ar []int, mx uint32) {
	if len(ar) < S {
		forSortI(ar)
		return
	}

	if mx < 2 { // 2..65536 goroutines
		mxI = 2
	} else if mx > 65536 {
		mxI = 65536
	} else {
		mxI = mx
	}
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
