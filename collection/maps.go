package collection

import "github.com/yuzhibiao/gfun/typex"

// Entries 返回 m 的键值对切片。map 无序，结果顺序不确定。
func Entries[K comparable, V any](m map[K]V) []typex.Pair[K, V] {
	r := make([]typex.Pair[K, V], 0, len(m))
	for k, v := range m {
		r = append(r, typex.Pair[K, V]{Key: k, Val: v})
	}
	return r
}

// Invert 交换键值。值不唯一时，保留哪个键由 map 遍历顺序决定（不确定）。
func Invert[K, V comparable](m map[K]V) map[V]K {
	r := make(map[V]K, len(m))
	for k, v := range m {
		r[v] = k
	}
	return r
}

// Merge 合并多个 map，靠后的覆盖靠前的。
func Merge[K comparable, V any](ms ...map[K]V) map[K]V {
	r := make(map[K]V)
	for _, m := range ms {
		for k, v := range m {
			r[k] = v
		}
	}
	return r
}

// FilterMap 返回满足 keep 的键值对组成的新 map。
// （不叫 Filter 是为了避免与 slice 的 Filter 重名。）
func FilterMap[K comparable, V any](m map[K]V, keep func(K, V) bool) map[K]V {
	r := make(map[K]V)
	for k, v := range m {
		if keep(k, v) {
			r[k] = v
		}
	}
	return r
}
