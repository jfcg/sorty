package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// uint array to be sorted
var arU []uint

// Checks if ar is sorted in ascending order.
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

var wgU sync.WaitGroup

// Concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
func SortU(ar []uint) {
	if len(ar) < S {
		forSortU(ar)
		return
	}
	arU = ar

	wgU.Add(1) // count self
	srtU(0, len(arU)-1)
	wgU.Wait()

	arU = nil
}

// assumes hi-lo >= S-1
func srtU(lo, hi int) {
	var l, h int
	var pv uint
start:
	pv = medianU(lo, hi)
	l, h = lo+1, hi-1 // medianU handles lo,hi positions

	for l <= h {
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

	if hi-l < S-1 { // hi range small?
		forSortU(arU[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortU(arU[lo : h+1])

			wgU.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortU(arU[lo : h+1])
	} else {
		wgU.Add(1)
		go srtU(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
