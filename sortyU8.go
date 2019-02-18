package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// uint64 array to be sorted
var ArU8 []uint64

func forSortU8(ar []uint64) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

func medianU8(l, h int) uint64 {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := ArU8[l], ArU8[m], ArU8[h]

	if vh < vl { // choose pivot as median of ArU8[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			ArU8[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			ArU8[m] = pv
		}

		ArU8[l], ArU8[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			ArU8[m] = pv
			ArU8[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			ArU8[m] = pv
			ArU8[l] = vl
		}
	}

	return pv
}

var wgU8 sync.WaitGroup

// Concurrently sorts ArU8 of type []uint64
func SortU8() {
	if len(ArU8) < S {
		forSortU8(ArU8)
		return
	}
	wgU8.Add(1) // count self
	srtU8(0, len(ArU8)-1)
	wgU8.Wait()
}

// assumes hi-lo >= S-1
func srtU8(lo, hi int) {
	var l, h int
	var pv uint64
start:
	pv = medianU8(lo, hi)
	l, h = lo+1, hi-1 // medianU8 handles lo,hi positions

	for l <= h {
		ct := false
		if ArU8[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if ArU8[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		ArU8[l], ArU8[h] = ArU8[h], ArU8[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		if hi > l {
			forSortU8(ArU8[l : hi+1])
		}

		if h-lo < S-1 { // lo range small?
			if h > lo {
				forSortU8(ArU8[lo : h+1])
			}
			wgU8.Done() // signal finish
			return    // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		if h > lo {
			forSortU8(ArU8[lo : h+1])
		}
	} else {
		wgU8.Add(1)
		go srtU8(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
