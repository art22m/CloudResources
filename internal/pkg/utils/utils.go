package utils

import "golang.org/x/exp/constraints"

func AtLeastOneNonZero(nums ...int) bool {
	for _, num := range nums {
		if num > 0 {
			return true
		}
	}
	return false
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}

	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}

	return b
}
