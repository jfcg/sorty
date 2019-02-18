package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// uint32 array to be sorted
var ArU4 []uint32

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
	vl, pv, vh := ArU4[l], ArU4[m], ArU4[h]

	if vh < vl { // choose pivot as median of ArU4[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			ArU4[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			ArU4[m] = pv
		}

		ArU4[l], ArU4[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			ArU4[m] = pv
			ArU4[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			ArU4[m] = pv
			ArU4[l] = vl
		}
	}

	return pv
}

var wgU4 sync.WaitGroup

// Concurrently sorts ArU4 of type []uint32
func SortU4() {
	if len(ArU4) < S {
		forSortU4(ArU4)
		return
	}
	wgU4.Add(1) // count self
	srtU4(0, len(ArU4)-1)
	wgU4.Wait()
}

// assumes hi-lo >= S-1
func srtU4(lo, hi int) {
	var l, h int
	var pv uint32
start:
	pv = medianU4(lo, hi)
	l, h = lo+1, hi-1 // medianU4 handles lo,hi positions

	for l <= h {
		ct := false
		if ArU4[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if ArU4[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		ArU4[l], ArU4[h] = ArU4[h], ArU4[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		if hi > l {
			forSortU4(ArU4[l : hi+1])
		}

		if h-lo < S-1 { // lo range small?
			if h > lo {
				forSortU4(ArU4[lo : h+1])
			}
			wgU4.Done() // signal finish
			return    // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		if h > lo {
			forSortU4(ArU4[lo : h+1])
		}
	} else {
		wgU4.Add(1)
		go srtU4(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
