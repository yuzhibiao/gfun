package typex

import "cmp"

// Clamp 将 v 限制在 [lo, hi] 区间内。
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Sum 求和，也适用于字符串拼接。
func Sum[T cmp.Ordered](s ...T) T {
	var r T
	for _, v := range s {
		r += v
	}
	return r
}
