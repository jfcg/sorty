package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// uintptr array to be sorted
var arP []uintptr

// Checks if ar is sorted in ascending order.
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

var wgP sync.WaitGroup

// Concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
func SortP(ar []uintptr) {
	if len(ar) < S {
		forSortP(ar)
		return
	}
	arP = ar

	wgP.Add(1) // count self
	srtP(0, len(arP)-1)
	wgP.Wait()

	arP = nil
}

// assumes hi-lo >= S-1
func srtP(lo, hi int) {
	var l, h int
	var pv uintptr
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

			wgP.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortP(arP[lo : h+1])
	} else {
		wgP.Add(1)
		go srtP(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
