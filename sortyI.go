/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "unsafe"

// IsSortedI returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedI(ar []int) int {
	s := unsafe.Sizeof(int(0))
	if s == 4 {
		return IsSortedI4(*(*[]int32)(unsafe.Pointer(&ar)))
	}
	if s == 8 {
		return IsSortedI8(*(*[]int64)(unsafe.Pointer(&ar)))
	}
	panic("sorty: IsSortedI: hw word size unknown")
}

// IsSortedU returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedU(ar []uint) int {
	s := unsafe.Sizeof(uint(0))
	if s == 4 {
		return IsSortedU4(*(*[]uint32)(unsafe.Pointer(&ar)))
	}
	if s == 8 {
		return IsSortedU8(*(*[]uint64)(unsafe.Pointer(&ar)))
	}
	panic("sorty: IsSortedU: hw word size unknown")
}

// IsSortedP returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedP(ar []uintptr) int {
	s := unsafe.Sizeof(uintptr(0))
	if s == 4 {
		return IsSortedU4(*(*[]uint32)(unsafe.Pointer(&ar)))
	}
	if s == 8 {
		return IsSortedU8(*(*[]uint64)(unsafe.Pointer(&ar)))
	}
	panic("sorty: IsSortedP: hw pointer size unknown")
}

// SortI concurrently sorts ar in ascending order.
func SortI(ar []int) {
	s := unsafe.Sizeof(int(0))
	if s == 4 {
		SortI4(*(*[]int32)(unsafe.Pointer(&ar)))
		return
	}
	if s == 8 {
		SortI8(*(*[]int64)(unsafe.Pointer(&ar)))
		return
	}
	panic("sorty: SortI: hw word size unknown")
}

// SortU concurrently sorts ar in ascending order.
func SortU(ar []uint) {
	s := unsafe.Sizeof(uint(0))
	if s == 4 {
		SortU4(*(*[]uint32)(unsafe.Pointer(&ar)))
		return
	}
	if s == 8 {
		SortU8(*(*[]uint64)(unsafe.Pointer(&ar)))
		return
	}
	panic("sorty: SortU: hw word size unknown")
}

// SortP concurrently sorts ar in ascending order.
func SortP(ar []uintptr) {
	s := unsafe.Sizeof(uintptr(0))
	if s == 4 {
		SortU4(*(*[]uint32)(unsafe.Pointer(&ar)))
		return
	}
	if s == 8 {
		SortU8(*(*[]uint64)(unsafe.Pointer(&ar)))
		return
	}
	panic("sorty: SortP: hw pointer size unknown")
}
