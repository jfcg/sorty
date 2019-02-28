package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// uint64 array to be sorted
var arU8 []uint64

// Checks if ar is sorted in ascending order.
func IsSortedU8(ar []uint64) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

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
	vl, pv, vh := arU8[l], arU8[m], arU8[h]

	if vh < vl { // choose pivot as median of arU8[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arU8[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arU8[m] = pv
		}

		arU8[l], arU8[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arU8[m] = pv
			arU8[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arU8[m] = pv
			arU8[l] = vl
		}
	}

	return pv
}

var wgU8 sync.WaitGroup

// Concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
func SortU8(ar []uint64) {
	if len(ar) < S {
		forSortU8(ar)
		return
	}
	arU8 = ar

	wgU8.Add(1) // count self
	srtU8(0, len(arU8)-1)
	wgU8.Wait()

	arU8 = nil
}

// assumes hi-lo >= S-1
func srtU8(lo, hi int) {
	var l, h int
	var pv uint64
start:
	pv = medianU8(lo, hi)
	l, h = lo+1, hi-1 // medianU8 handles lo,hi positions

	for l <= h {
		swap := true
		if arU8[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arU8[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arU8[l], arU8[h] = arU8[h], arU8[l]
			h--
			l++
		}
	}

	if hi-l < S-1 { // hi range small?
		forSortU8(arU8[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortU8(arU8[lo : h+1])

			wgU8.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortU8(arU8[lo : h+1])
	} else {
		wgU8.Add(1)
		go srtU8(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
