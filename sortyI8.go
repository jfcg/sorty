package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// int64 array to be sorted
var arI8 []int64

func forSortI8(ar []int64) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

func medianI8(l, h int) int64 {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arI8[l], arI8[m], arI8[h]

	if vh < vl { // choose pivot as median of arI8[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arI8[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arI8[m] = pv
		}

		arI8[l], arI8[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arI8[m] = pv
			arI8[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arI8[m] = pv
			arI8[l] = vl
		}
	}

	return pv
}

var wgI8 sync.WaitGroup

// Concurrently sorts ar. Should not be called by multiple goroutines at the same time.
func SortI8(ar []int64) {
	if len(ar) < S {
		forSortI8(ar)
		return
	}
	arI8 = ar

	wgI8.Add(1) // count self
	srtI8(0, len(arI8)-1)
	wgI8.Wait()

	arI8 = nil
}

// assumes hi-lo >= S-1
func srtI8(lo, hi int) {
	var l, h int
	var pv int64
start:
	pv = medianI8(lo, hi)
	l, h = lo+1, hi-1 // medianI8 handles lo,hi positions

	for l <= h {
		ct := false
		if arI8[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if arI8[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		arI8[l], arI8[h] = arI8[h], arI8[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		forSortI8(arI8[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortI8(arI8[lo : h+1])

			wgI8.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortI8(arI8[lo : h+1])
	} else {
		wgI8.Add(1)
		go srtI8(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
