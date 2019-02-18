package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// uintptr array to be sorted
var ArP []uintptr

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
	vl, pv, vh := ArP[l], ArP[m], ArP[h]

	if vh < vl { // choose pivot as median of ArP[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			ArP[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			ArP[m] = pv
		}

		ArP[l], ArP[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			ArP[m] = pv
			ArP[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			ArP[m] = pv
			ArP[l] = vl
		}
	}

	return pv
}

var wgP sync.WaitGroup

// Concurrently sorts ArP of type []uintptr
func SortP() {
	if len(ArP) < S {
		forSortP(ArP)
		return
	}
	wgP.Add(1) // count self
	srtP(0, len(ArP)-1)
	wgP.Wait()
}

// assumes hi-lo >= S-1
func srtP(lo, hi int) {
	var l, h int
	var pv uintptr
start:
	pv = medianP(lo, hi)
	l, h = lo+1, hi-1 // medianP handles lo,hi positions

	for l <= h {
		ct := false
		if ArP[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if ArP[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		ArP[l], ArP[h] = ArP[h], ArP[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		if hi > l {
			forSortP(ArP[l : hi+1])
		}

		if h-lo < S-1 { // lo range small?
			if h > lo {
				forSortP(ArP[lo : h+1])
			}
			wgP.Done() // signal finish
			return    // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		if h > lo {
			forSortP(ArP[lo : h+1])
		}
	} else {
		wgP.Add(1)
		go srtP(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
