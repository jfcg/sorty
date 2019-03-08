package sorty

// Concurrent Sorting
// Author: Serhat Şevki Dinçer, jfcgaussATgmail

import "sync/atomic"

// uint64 array to be sorted
var arU8 []uint64

// IsSortedU8 checks if ar is sorted in ascending order.
func IsSortedU8(ar []uint64) bool {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return false
		}
	}
	return true
}

// insertion sort
func insertionU8(ar []uint64) {
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
func ipU8(pv, vl, vh uint64) (a, b, c uint64, r int) {
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
func medianU8(l, h int) uint64 {
	// lo, med, hi
	m := mean(l, h)
	vl, pv, vh := arU8[l], arU8[m], arU8[h]

	// intermediates
	a, b := mean(l, m), mean(m, h)
	va, vb := arU8[a], arU8[b]

	// put lo, med, hi in order
	if vh < vl {
		vl, vh = vh, vl
	}
	vl, pv, vh, _ = ipU8(pv, vl, vh)

	// update pivot with intermediates
	if vb < va {
		va, vb = vb, va
	}
	va, pv, vb, r := ipU8(pv, va, vb)

	// if pivot was out of [va, vb]
	if r == 1 {
		vl, va, pv, _ = ipU8(vl, va, pv)
	} else if r == -1 {
		pv, vb, vh, _ = ipU8(vh, pv, vb)
	}

	// here: vl <= va <= pv <= vb <= vh
	arU8[l], arU8[m], arU8[h] = vl, pv, vh
	arU8[a], arU8[b] = va, vb
	return pv
}

var ngU8, mxU8 uint32 // number of sorting goroutines, max limit
var doneU8 = make(chan bool, 1)

// SortU8 concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortU8(ar []uint64, mx uint32) {
	if len(ar) <= Mli {
		insertionU8(ar)
		return
	}

	mxU8 = sat(mx)
	arU8 = ar

	ngU8 = 1 // count self
	gsrtU8(0, len(arU8)-1)
	<-doneU8

	arU8 = nil
}

func gsrtU8(lo, hi int) {
	srtU8(lo, hi)

	if atomic.AddUint32(&ngU8, ^uint32(0)) == 0 { // decrease goroutine counter
		doneU8 <- false // we are the last, all done
	}
}

// assumes hi-lo >= Mli
func srtU8(lo, hi int) {
	var l, h int
start:
	l, h = lo+1, hi-1 // medianU8 handles lo,hi positions

	for pv := medianU8(lo, hi); l <= h; {
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

	if h-lo < hi-l {
		h, hi = hi, h // [lo,h] is the bigger range
		l, lo = lo, l
	}

	if hi-l >= Mli { // two big ranges?

		// max goroutines? range not big enough for new goroutine?
		// not atomic but good enough
		if ngU8 >= mxU8 || hi-l <= 2*Mli {
			srtU8(l, hi) // start a recursive (slave) sort on the smaller range
			hi = h
			goto start
		}

		atomic.AddUint32(&ngU8, 1) // increase goroutine counter
		go gsrtU8(lo, h)           // start a goroutine on the bigger range
		lo = l
		goto start
	}

	insertionU8(arU8[l : hi+1])

	if h-lo < Mli { // two small ranges?
		insertionU8(arU8[lo : h+1])
		return
	}

	hi = h
	goto start
}
