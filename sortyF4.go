package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// float32 array to be sorted
var ArF4 []float32

func forSortF4(ar []float32) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

func medianF4(l, h int) float32 {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := ArF4[l], ArF4[m], ArF4[h]

	if vh < vl { // choose pivot as median of ArF4[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			ArF4[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			ArF4[m] = pv
		}

		ArF4[l], ArF4[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			ArF4[m] = pv
			ArF4[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			ArF4[m] = pv
			ArF4[l] = vl
		}
	}

	return pv
}

var wgF4 sync.WaitGroup

// Concurrently sorts ArF4 of type []float32
func SortF4() {
	if len(ArF4) < S {
		forSortF4(ArF4)
		return
	}
	wgF4.Add(1) // count self
	srtF4(0, len(ArF4)-1)
	wgF4.Wait()
}

// assumes hi-lo >= S-1
func srtF4(lo, hi int) {
	var l, h int
	var pv float32
start:
	pv = medianF4(lo, hi)
	l, h = lo+1, hi-1 // medianF4 handles lo,hi positions

	for l <= h {
		ct := false
		if ArF4[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if ArF4[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		ArF4[l], ArF4[h] = ArF4[h], ArF4[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		forSortF4(ArF4[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortF4(ArF4[lo : h+1])

			wgF4.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortF4(ArF4[lo : h+1])
	} else {
		wgF4.Add(1)
		go srtF4(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
