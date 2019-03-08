package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// uint32 array to be sorted
var arU4 []uint32

// IsSortedU4 checks if ar is sorted in ascending order.
func IsSortedU4(ar []uint32) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

// insertion sort
func insertionU4(ar []uint32) {
	for h := 1; h < len(ar); h++ {
		v, l := ar[h], h-1
		if v < ar[l] {
			for {
				ar[l+1] = ar[l]
				l--
				if l < 0 || v >= ar[l] {
					break
				}
			}
			ar[l+1] = v
		}
	}
}

// given vl <= vh, inserts pv in the middle
// returns vl <= pv <= vh
func ipU4(pv, vl, vh uint32) (a, b, c uint32, r int) {
	if pv > vh {
		vh, pv = pv, vh
		r = 1
	} else if pv < vl {
		vl, pv = pv, vl
		r = -1
	}
	return vl, pv, vh, r
}

// return pivot as median of five scattered values
func medianU4(l, h int) uint32 {
	// lo, med, hi
	m := mean(l, h)
	vl, pv, vh := arU4[l], arU4[m], arU4[h]

	// intermediates
	a, b := mean(l, m), mean(m, h)
	va, vb := arU4[a], arU4[b]

	// put lo, med, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipU4(pv, vl, vh)

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipU4(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipU4(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipU4(vh, pv, vb)
	}

	// here: vl <= va <= pv <= vb <= vh
	arU4[l], arU4[m], arU4[h] = vl, pv, vh
	arU4[a], arU4[b] = va, vb
	return pv
}

var ngU4, mxU4 uint32 // number of sorting goroutines, max limit
var doneU4 = make(chan bool, 1)

// SortU4 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortU4(ar []uint32, mx uint32) {
	if len(ar) <= Mli {
		insertionU4(ar)
		return
	}

	mxU4 = sat(mx)
	arU4 = ar

	ngU4 = 1 // count self
	gsrtU4(0, len(arU4)-1)
	<-doneU4

	arU4 = nil
}

func gsrtU4(lo, hi int) {
	srtU4(lo, hi)

	if atomic.AddUint32(&ngU4, ^uint32(0)) == 0 { // decrease goroutine counter
		doneU4 <- false // we are the last, all done
	}
}

// assumes hi-lo >= Mli
func srtU4(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianU4 handles lo,hi positions

	for pv := medianU4(lo, hi); l <= h; {
		swap := true
		if arU4[h] >= pv { // extend ranges in balance
			h--
			swap = false
		}
		if arU4[l] <= pv {
			l++
			swap = false
		}

		if swap {
			arU4[l], arU4[h] = arU4[h], arU4[l]
			h--
			l++
		}
	}

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= Mli { // two big ranges?

		// max goroutines? range not big enough for new goroutine?
		// not atomic but good enough
		if ngU4 >= mxU4 || hi-l <= 2*Mli {
			srtU4(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngU4, 1) // increase goroutine counter
		go gsrtU4(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	insertionU4(arU4[l : hi+1])

	if h-lo < Mli { // two small ranges?
		insertionU4(arU4[lo : h+1])
		return
	}

	hi = h
	goto start
}
