package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// float32 array to be sorted
var arF4 []float32

// Checks if ar is sorted in ascending order.
func IsSortedF4(ar []float32) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

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
	vl, pv, vh := arF4[l], arF4[m], arF4[h]

	if vh < vl { // choose pivot as median of arF4[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arF4[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arF4[m] = pv
		}

		arF4[l], arF4[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arF4[m] = pv
			arF4[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arF4[m] = pv
			arF4[l] = vl
		}
	}

	return pv
}

var wgF4 sync.WaitGroup

// Concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
func SortF4(ar []float32) {
	if len(ar) < S {
		forSortF4(ar)
		return
	}
	arF4 = ar

	wgF4.Add(1) // count self
	srtF4(0, len(arF4)-1)
	wgF4.Wait()

	arF4 = nil
}

// assumes hi-lo >= S-1
func srtF4(lo, hi int) {
	var l, h int
	var pv float32
start:
	pv = medianF4(lo, hi)
	l, h = lo+1, hi-1 // medianF4 handles lo,hi positions

	for l <= h {
		swap := true
		if arF4[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arF4[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arF4[l], arF4[h] = arF4[h], arF4[l]
			h--
			l++
		}
	}

	if hi-l < S-1 { // hi range small?
		forSortF4(arF4[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortF4(arF4[lo : h+1])

			wgF4.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortF4(arF4[lo : h+1])
	} else {
		wgF4.Add(1)
		go srtF4(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
