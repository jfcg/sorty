package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// int32 array to be sorted
var ArI4 []int32

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
	vl, pv, vh := ArI4[l], ArI4[m], ArI4[h]

	if vh < vl { // choose pivot as median of ArI4[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			ArI4[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			ArI4[m] = pv
		}

		ArI4[l], ArI4[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			ArI4[m] = pv
			ArI4[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			ArI4[m] = pv
			ArI4[l] = vl
		}
	}

	return pv
}

var wgI4 sync.WaitGroup

// Concurrently sorts ArI4 of type []int32
func SortI4() {
	if len(ArI4) < S {
		forSortI4(ArI4)
		return
	}
	wgI4.Add(1) // count self
	srtI4(0, len(ArI4)-1)
	wgI4.Wait()
}

// assumes hi-lo >= S-1
func srtI4(lo, hi int) {
	var l, h int
	var pv int32
start:
	pv = medianI4(lo, hi)
	l, h = lo+1, hi-1 // medianI4 handles lo,hi positions

	for l <= h {
		ct := false
		if ArI4[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if ArI4[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		ArI4[l], ArI4[h] = ArI4[h], ArI4[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		forSortI4(ArI4[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortI4(ArI4[lo : h+1])

			wgI4.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortI4(ArI4[lo : h+1])
	} else {
		wgI4.Add(1)
		go srtI4(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
