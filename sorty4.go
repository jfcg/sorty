// +build 386 arm mips mipsle

package sorty

import "unsafe"

func init() {
	if unsafe.Sizeof(int(0)) != 4 || unsafe.Sizeof(uint(0)) != 4 || unsafe.Sizeof(uintptr(0)) != 4 {
		panic("Architecture word/pointer size mismatch")
	}
}

// IsSortedI checks if ar is sorted in ascending order.
func IsSortedI(ar []int) bool {
	ap := *(*[]int32)(unsafe.Pointer(&ar))
	return IsSortedI4(ap)
}

// IsSortedU checks if ar is sorted in ascending order.
func IsSortedU(ar []uint) bool {
	ap := *(*[]uint32)(unsafe.Pointer(&ar))
	return IsSortedU4(ap)
}

// IsSortedP checks if ar is sorted in ascending order.
func IsSortedP(ar []uintptr) bool {
	ap := *(*[]uint32)(unsafe.Pointer(&ar))
	return IsSortedU4(ap)
}

// SortI concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
// SortI calls SortI4 (on 32-bit) or SortI8 (on 64-bit).
func SortI(ar []int, mx uint32) {
	ap := *(*[]int32)(unsafe.Pointer(&ar))
	SortI4(ap, mx)
}

// SortU concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
// SortU calls SortU4 (on 32-bit) or SortU8 (on 64-bit).
func SortU(ar []uint, mx uint32) {
	ap := *(*[]uint32)(unsafe.Pointer(&ar))
	SortU4(ap, mx)
}

// SortP concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
// SortP calls SortU4 (on 32-bit) or SortU8 (on 64-bit).
func SortP(ar []uintptr, mx uint32) {
	ap := *(*[]uint32)(unsafe.Pointer(&ar))
	SortU4(ap, mx)
}
