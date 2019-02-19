package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// int array to be sorted
var ArI []int

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
	vl, pv, vh := ArI[l], ArI[m], ArI[h]

	if vh < vl { // choose pivot as median of ArI[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			ArI[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			ArI[m] = pv
		}

		ArI[l], ArI[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			ArI[m] = pv
			ArI[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			ArI[m] = pv
			ArI[l] = vl
		}
	}

	return pv
}

var wgI sync.WaitGroup

// Concurrently sorts ArI of type []int
func SortI() {
	if len(ArI) < S {
		forSortI(ArI)
		return
	}
	wgI.Add(1) // count self
	srtI(0, len(ArI)-1)
	wgI.Wait()
}

// assumes hi-lo >= S-1
func srtI(lo, hi int) {
	var l, h int
	var pv int
start:
	pv = medianI(lo, hi)
	l, h = lo+1, hi-1 // medianI handles lo,hi positions

	for l <= h {
		ct := false
		if ArI[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if ArI[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		ArI[l], ArI[h] = ArI[h], ArI[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		forSortI(ArI[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortI(ArI[lo : h+1])

			wgI.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortI(ArI[lo : h+1])
	} else {
		wgI.Add(1)
		go srtI(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
