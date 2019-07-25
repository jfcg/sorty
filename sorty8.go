// +build amd64 arm64 ppc64 ppc64le mips64 mips64le s390x

package sorty

import "unsafe"

func init() {
	if unsafe.Sizeof(int(0)) != 8 || unsafe.Sizeof(uint(0)) != 8 || unsafe.Sizeof(uintptr(0)) != 8 {
		panic("sorty: architecture word/pointer size mismatch")
	}
}

// IsSortedI checks if ar is sorted in ascending order.
func IsSortedI(ar []int) bool {
	return IsSortedI8(*(*[]int64)(unsafe.Pointer(&ar)))
}

// IsSortedU checks if ar is sorted in ascending order.
func IsSortedU(ar []uint) bool {
	return IsSortedU8(*(*[]uint64)(unsafe.Pointer(&ar)))
}

// IsSortedP checks if ar is sorted in ascending order.
func IsSortedP(ar []uintptr) bool {
	return IsSortedU8(*(*[]uint64)(unsafe.Pointer(&ar)))
}

// SortI concurrently sorts ar in ascending order. mx is the maximum number
// of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortI(ar []int, mx uint32) {
	SortI8(*(*[]int64)(unsafe.Pointer(&ar)), mx)
}

// SortU concurrently sorts ar in ascending order. mx is the maximum number
// of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortU(ar []uint, mx uint32) {
	SortU8(*(*[]uint64)(unsafe.Pointer(&ar)), mx)
}

// SortP concurrently sorts ar in ascending order. mx is the maximum number
// of goroutines used for sorting simultaneously, saturated to [2, 65535].
func SortP(ar []uintptr, mx uint32) {
	SortU8(*(*[]uint64)(unsafe.Pointer(&ar)), mx)
}
