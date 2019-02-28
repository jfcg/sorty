package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// float64 array to be sorted
var arF8 []float64

// Checks if ar is sorted in ascending order.
func IsSortedF8(ar []float64) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

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
	vl, pv, vh := arF8[l], arF8[m], arF8[h]

	if vh < vl { // choose pivot as median of arF8[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arF8[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arF8[m] = pv
		}

		arF8[l], arF8[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arF8[m] = pv
			arF8[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arF8[m] = pv
			arF8[l] = vl
		}
	}

	return pv
}

var wgF8 sync.WaitGroup

// Concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
func SortF8(ar []float64) {
	if len(ar) < S {
		forSortF8(ar)
		return
	}
	arF8 = ar

	wgF8.Add(1) // count self
	srtF8(0, len(arF8)-1)
	wgF8.Wait()

	arF8 = nil
}

// assumes hi-lo >= S-1
func srtF8(lo, hi int) {
	var l, h int
	var pv float64
start:
	pv = medianF8(lo, hi)
	l, h = lo+1, hi-1 // medianF8 handles lo,hi positions

	for l <= h {
		swap := true
		if arF8[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arF8[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arF8[l], arF8[h] = arF8[h], arF8[l]
			h--
			l++
		}
	}

	if hi-l < S-1 { // hi range small?
		forSortF8(arF8[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortF8(arF8[lo : h+1])

			wgF8.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortF8(arF8[lo : h+1])
	} else {
		wgF8.Add(1)
		go srtF8(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
