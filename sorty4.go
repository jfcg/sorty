// +build 386 arm mips mipsle wasm

/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import "unsafe"

func init() {
	if unsafe.Sizeof(int(0)) != 4 || unsafe.Sizeof(uint(0)) != 4 ||
		unsafe.Sizeof(uintptr(0)) != 4 {
		panic("sorty: architecture word/pointer size mismatch")
	}
}

// IsSortedI returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedI(ar []int) int {
	return IsSortedI4(*(*[]int32)(unsafe.Pointer(&ar)))
}

// IsSortedU returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedU(ar []uint) int {
	return IsSortedU4(*(*[]uint32)(unsafe.Pointer(&ar)))
}

// IsSortedP returns 0 if ar is sorted in ascending order,
// otherwise it returns i > 0 with ar[i] < ar[i-1]
func IsSortedP(ar []uintptr) int {
	return IsSortedU4(*(*[]uint32)(unsafe.Pointer(&ar)))
}

// SortI concurrently sorts ar in ascending order.
func SortI(ar []int) {
	SortI4(*(*[]int32)(unsafe.Pointer(&ar)))
}

// SortU concurrently sorts ar in ascending order.
func SortU(ar []uint) {
	SortU4(*(*[]uint32)(unsafe.Pointer(&ar)))
}

// SortP concurrently sorts ar in ascending order.
func SortP(ar []uintptr) {
	SortU4(*(*[]uint32)(unsafe.Pointer(&ar)))
}
