package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// uint array to be sorted
var arU []uint

// IsSortedU checks if ar is sorted in ascending order.
func IsSortedU(ar []uint) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

func forSortU(ar []uint) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

func medianU(l, h int) uint {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arU[l], arU[m], arU[h]

	if vh < vl { // choose pivot as median of arU[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arU[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arU[m] = pv
		}

		arU[l], arU[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arU[m] = pv
			arU[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arU[m] = pv
			arU[l] = vl
		}
	}

	return pv
}

var ngU, mxU uint32 // number of sorting goroutines, max limit
var doneU = make(chan bool, 1)

// SortU concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting, saturated to [2, 65536].
func SortU(ar []uint, mx uint32) {
	if len(ar) < S {
		forSortU(ar)
		return
	}

	if mx < 2 { // 2..65536 goroutines
		mxU = 2
	} else if mx > 65536 {
		mxU = 65536
	} else {
		mxU = mx
	}
	arU = ar

	ngU = 1 // count self
	gsrtU(0, len(arU)-1)
	<-doneU

	arU = nil
}

func gsrtU(lo, hi int) {
	srtU(lo, hi)

	if atomic.AddUint32(&ngU, ^uint32(0)) == 0 { // decrease goroutine counter
		doneU <- false // we are the last, all done
	}
}

// assumes hi-lo >= S-1
func srtU(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianU handles lo,hi positions

	for pv := medianU(lo, hi); l <= h; {
		swap := true
		if arU[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arU[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arU[l], arU[h] = arU[h], arU[l]
			h--
			l++
		}
	}

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= S-1 { // two big ranges?

		if ngU >= mxU { // max number of goroutines? not atomic but good enough
			srtU(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngU, 1) // increase goroutine counter
		go gsrtU(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	forSortU(arU[l : hi+1])

	if h-lo < S-1 { // two small ranges?
		forSortU(arU[lo : h+1])
		return
	}

	hi = h
	goto start
}
