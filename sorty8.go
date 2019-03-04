// +build amd64 arm64 ppc64 ppc64le mips64 mips64le s390x

package sorty

import "unsafe"

func init() {
	if unsafe.Sizeof(int(0)) != 8 || unsafe.Sizeof(uint(0)) != 8 || unsafe.Sizeof(uintptr(0)) != 8 {
		panic("Architecture word/pointer size mismatch")
	}
}

// IsSortedI checks if ar is sorted in ascending order.
func IsSortedI(ar []int) bool {
	ap := *(*[]int64)(unsafe.Pointer(&ar))
	return IsSortedI8(ap)
}

// IsSortedU checks if ar is sorted in ascending order.
func IsSortedU(ar []uint) bool {
	ap := *(*[]uint64)(unsafe.Pointer(&ar))
	return IsSortedU8(ap)
}

// IsSortedP checks if ar is sorted in ascending order.
func IsSortedP(ar []uintptr) bool {
	ap := *(*[]uint64)(unsafe.Pointer(&ar))
	return IsSortedU8(ap)
}

// SortI concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
// SortI calls SortI4 (on 32-bit) or SortI8 (on 64-bit).
func SortI(ar []int, mx uint32) {
	ap := *(*[]int64)(unsafe.Pointer(&ar))
	SortI8(ap, mx)
}

// SortU concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
// SortU calls SortU4 (on 32-bit) or SortU8 (on 64-bit).
func SortU(ar []uint, mx uint32) {
	ap := *(*[]uint64)(unsafe.Pointer(&ar))
	SortU8(ap, mx)
}

// SortP concurrently sorts ar in ascending order. Should not be called by multiple goroutines at the same time.
// mx is the maximum number of goroutines used for sorting simultaneously, saturated to [2, 65535].
// SortP calls SortU4 (on 32-bit) or SortU8 (on 64-bit).
func SortP(ar []uintptr, mx uint32) {
	ap := *(*[]uint64)(unsafe.Pointer(&ar))
	SortU8(ap, mx)
}
