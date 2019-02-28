package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// int array to be sorted
var arI []int

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

var wgI sync.WaitGroup

// Concurrently sorts ar. Should not be called by multiple goroutines at the same time.
func SortI(ar []int) {
	if len(ar) < S {
		forSortI(ar)
		return
	}
	arI = ar

	wgI.Add(1) // count self
	srtI(0, len(arI)-1)
	wgI.Wait()

	arI = nil
}

// assumes hi-lo >= S-1
func srtI(lo, hi int) {
	var l, h int
	var pv int
start:
	pv = medianI(lo, hi)
	l, h = lo+1, hi-1 // medianI handles lo,hi positions

	for l <= h {
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

	if hi-l < S-1 { // hi range small?
		forSortI(arI[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortI(arI[lo : h+1])

			wgI.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortI(arI[lo : h+1])
	} else {
		wgI.Add(1)
		go srtI(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
