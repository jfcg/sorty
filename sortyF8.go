package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// float64 array to be sorted
var ArF8 []float64

func forSortF8(ar []float64) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

func medianF8(l, h int) float64 {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := ArF8[l], ArF8[m], ArF8[h]

	if vh < vl { // choose pivot as median of ArF8[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			ArF8[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			ArF8[m] = pv
		}

		ArF8[l], ArF8[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			ArF8[m] = pv
			ArF8[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			ArF8[m] = pv
			ArF8[l] = vl
		}
	}

	return pv
}

var wgF8 sync.WaitGroup

// Concurrently sorts ArF8 of type []float64
func SortF8() {
	if len(ArF8) < S {
		forSortF8(ArF8)
		return
	}
	wgF8.Add(1) // count self
	srtF8(0, len(ArF8)-1)
	wgF8.Wait()
}

// assumes hi-lo >= S-1
func srtF8(lo, hi int) {
	var l, h int
	var pv float64
start:
	pv = medianF8(lo, hi)
	l, h = lo+1, hi-1 // medianF8 handles lo,hi positions

	for l <= h {
		ct := false
		if ArF8[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if ArF8[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		ArF8[l], ArF8[h] = ArF8[h], ArF8[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		forSortF8(ArF8[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortF8(ArF8[lo : h+1])

			wgF8.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortF8(ArF8[lo : h+1])
	} else {
		wgF8.Add(1)
		go srtF8(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
