package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync"

// Min array size for concurrent Sort()
const S = 24

// Array to be sorted
var Ar []uint64

func forSort(ar []uint64) {
	for h := len(ar) - 1; h > 0; h-- {
		for l := h - 1; l >= 0; l-- {
			if ar[h] < ar[l] {
				ar[l], ar[h] = ar[h], ar[l]
			}
		}
	}
}

func median(l, h int) (int, int, uint64) {
	m := int(uint(l+h) >> 1) // avoid overflow
	vl, pv, vh := Ar[l], Ar[m], Ar[h]

	if vh < vl { // choose pivot as median of Ar[l,m,h]
		vl, vh = vh, vl

		if pv > vh {
			vh, pv = pv, vh
			Ar[m] = pv
		} else if pv < vl {
			vl, pv = pv, vl
			Ar[m] = pv
		}

		Ar[l], Ar[h] = vl, vh
		h-- // update indices
		l++
	} else {
		if pv > vh {
			vh, pv = pv, vh
			Ar[m] = pv
			Ar[h] = vh
			h--
		} else if pv < vl {
			vl, pv = pv, vl
			Ar[m] = pv
			Ar[l] = vl
			l++
		}
	}

	return l, h, pv
}

var wg sync.WaitGroup

// Concurrently sorts Ar[]
func Sort() {
	if len(Ar) < S {
		forSort(Ar)
		return
	}
	wg.Add(1) // count self
	srt(0, len(Ar)-1)
	wg.Wait()
}

// assumes hi-lo >= S-1
func srt(lo, hi int) {
	var l, h int
	var pv uint64
start:
	l, h, pv = median(lo, hi)

	for l <= h {
		ct := false
		if Ar[h] >= pv { // extend ranges in balance
			h--
			ct = true
		}
		if Ar[l] <= pv {
			l++
			ct = true
		}
		if ct {
			continue
		}

		Ar[l], Ar[h] = Ar[h], Ar[l]
		h--
		l++
	}

	if hi-l < S-1 { // hi range small?
		if hi > l {
			forSort(Ar[l : hi+1])
		}

		if h-lo < S-1 { // lo range small?
			if h > lo {
				forSort(Ar[lo : h+1])
			}
			wg.Done() // signal finish
			return    // done with two small ranges
		}

		hi = h // continue with big lo range
		goto start
	}

	if h-lo < S-1 { // lo range small?
		if h > lo {
			forSort(Ar[lo : h+1])
		}
	} else {
		wg.Add(1)
		go srt(lo, h) // two big ranges, handle big lo range in another goroutine
	}

	lo = l // continue with big hi range
	goto start
}
