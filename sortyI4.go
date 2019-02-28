package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// int32 array to be sorted
var arI4 []int32

// Checks if ar is sorted in ascending order.
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

var wgI4 sync.WaitGroup

// Concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
func SortI4(ar []int32) {
	if len(ar) < S {
		forSortI4(ar)
		return
	}
	arI4 = ar

	wgI4.Add(1) // count self
	srtI4(0, len(arI4)-1)
	wgI4.Wait()

	arI4 = nil
}

// assumes hi-lo >= S-1
func srtI4(lo, hi int) {
	var l, h int
	var pv int32
start:
	pv = medianI4(lo, hi)
	l, h = lo+1, hi-1 // medianI4 handles lo,hi positions

	for l <= h {
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

	if hi-l < S-1 { // hi range small?
		forSortI4(arI4[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortI4(arI4[lo : h+1])

			wgI4.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortI4(arI4[lo : h+1])
	} else {
		wgI4.Add(1)
		go srtI4(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
