// Package collection 提供切片、map、集合的泛型操作。
//
// 标准库 slices/maps 已覆盖的能力（排序、查找、克隆等）不在本包重复实现。
package collection

// Map 对 s 中每个元素应用 fn，返回新切片。
func Map[T, U any](s []T, fn func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = fn(v)
	}
	return r
}

// Filter 返回 s 中满足 keep 的元素，保持原顺序。
func Filter[T any](s []T, keep func(T) bool) []T {
	r := make([]T, 0, len(s))
	for _, v := range s {
		if keep(v) {
			r = append(r, v)
		}
	}
	return r
}

// Reduce 以 init 为初始值，从左到右折叠 s。
func Reduce[T, U any](s []T, init U, fn func(U, T) U) U {
	acc := init
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// Chunk 将 s 按 size 分段，最后一段可能不足 size。size <= 0 时返回 nil。
func Chunk[T any](s []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	var r [][]T
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		r = append(r, s[i:end:end])
	}
	return r
}

// Unique 返回去重后的元素，保留首次出现的顺序。
func Unique[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	r := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			r = append(r, v)
		}
	}
	return r
}

// GroupBy 按 key 函数把元素分组。
func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	r := make(map[K][]T)
	for _, v := range s {
		k := key(v)
		r[k] = append(r[k], v)
	}
	return r
}
