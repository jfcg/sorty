package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// string array to be sorted
var arS []string

func forSortS(ar []string) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

func medianS(l, h int) string {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := arS[l], arS[m], arS[h]

	if vh < vl { // choose pivot as median of arS[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			arS[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			arS[m] = pv
		}

		arS[l], arS[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			arS[m] = pv
			arS[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			arS[m] = pv
			arS[l] = vl
		}
	}

	return pv
}

var wgS sync.WaitGroup

// Concurrently sorts ar. Should not be called by multiple goroutines at the same time.
func SortS(ar []string) {
	if len(ar) < S {
		forSortS(ar)
		return
	}
	arS = ar

	wgS.Add(1) // count self
	srtS(0, len(arS)-1)
	wgS.Wait()

	arS = nil
}

// assumes hi-lo >= S-1
func srtS(lo, hi int) {
	var l, h int
	var pv string
start:
	pv = medianS(lo, hi)
	l, h = lo+1, hi-1 // medianS handles lo,hi positions

	for l <= h {
		ct := false
		if arS[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if arS[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		arS[l], arS[h] = arS[h], arS[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		forSortS(arS[l : hi+1])

		if h-lo < S-1 { // lo range small?
			forSortS(arS[lo : h+1])

			wgS.Done() // signal finish
			return      // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		forSortS(arS[lo : h+1])
	} else {
		wgS.Add(1)
		go srtS(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
