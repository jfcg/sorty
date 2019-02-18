package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// string array to be sorted
var ArS []string

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
	vl, pv, vh := ArS[l], ArS[m], ArS[h]

	if vh < vl { // choose pivot as median of ArS[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			ArS[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			ArS[m] = pv
		}

		ArS[l], ArS[h] = vl, vh
	} else {
		if pv > vh {
			vh, pv = pv, vh
			ArS[m] = pv
			ArS[h] = vh
		} else if pv < vl {
			vl, pv = pv, vl
			ArS[m] = pv
			ArS[l] = vl
		}
	}

	return pv
}

var wgS sync.WaitGroup

// Concurrently sorts ArS of type []string
func SortS() {
	if len(ArS) < S {
		forSortS(ArS)
		return
	}
	wgS.Add(1) // count self
	srtS(0, len(ArS)-1)
	wgS.Wait()
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
		if ArS[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if ArS[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		ArS[l], ArS[h] = ArS[h], ArS[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		if hi > l {
			forSortS(ArS[l : hi+1])
		}

		if h-lo < S-1 { // lo range small?
			if h > lo {
				forSortS(ArS[lo : h+1])
			}
			wgS.Done() // signal finish
			return    // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		if h > lo {
			forSortS(ArS[lo : h+1])
		}
	} else {
		wgS.Add(1)
		go srtS(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
