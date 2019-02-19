package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// uint array to be sorted
var ArU []uint

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
	vl, pv, vh := ArU[l], ArU[m], ArU[h]

	if vh < vl { // choose pivot as median of ArU[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			ArU[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			ArU[m] = pv
		}

		ArU[l], ArU[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			ArU[m] = pv
			ArU[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			ArU[m] = pv
			ArU[l] = vl
		}
	}

	return pv
}

var wgU sync.WaitGroup

// Concurrently sorts ArU of type []uint
func SortU() {
	if len(ArU) < S {
		forSortU(ArU)
		return
	}
	wgU.Add(1) // count self
	srtU(0, len(ArU)-1)
	wgU.Wait()
}

// assumes hi-lo >= S-1
func srtU(lo, hi int) {
	var l, h int
	var pv uint
start:
	pv = medianU(lo, hi)
	l, h = lo+1, hi-1 // medianU handles lo,hi positions

	for l <= h {
		ct := false
		if ArU[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if ArU[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		ArU[l], ArU[h] = ArU[h], ArU[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		forSortU(ArU[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortU(ArU[lo : h+1])

			wgU.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortU(ArU[lo : h+1])
	} else {
		wgU.Add(1)
		go srtU(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
