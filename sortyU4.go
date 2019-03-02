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

var ngU4, mxU4 uint32 // number of sorting goroutines, max limit
var doneU4 = make(chan bool, 1)

// SortU4 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting, saturated to [2, 65536].
func SortU4(ar []uint32, mx uint32) {
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
	gsrtU4(0, len(arU4)-1)
	<-doneU4

	arU4 = nil
}

func gsrtU4(lo, hi int) {
	srtU4(lo, hi)

	if atomic.AddUint32(&ngU4, ^uint32(0)) == 0 { // decrease goroutine counter
		doneU4 <- false // we are the last, all done
	}
}

// assumes hi-lo >= S-1
func srtU4(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianU4 handles lo,hi positions

	for pv := medianU4(lo, hi); l <= h; {
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

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= S-1 { // two big ranges?

		if ngU4 >= mxU4 { // max number of goroutines? not atomic but good enough
			srtU4(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngU4, 1) // increase goroutine counter
		go gsrtU4(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	forSortU4(arU4[l : hi+1])

	if h-lo < S-1 { // two small ranges?
		forSortU4(arU4[lo : h+1])
		return
	}

	hi = h
	goto start
}
