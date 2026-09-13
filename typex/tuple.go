package typex

// Pair 是二元组。
type Pair[K, V any] struct {
	Key K
	Val V
}

// P 是构造 Pair 的简写。
func P[K, V any](k K, v V) Pair[K, V] {
	return Pair[K, V]{Key: k, Val: v}
}

// Zip 将两个切片按位置配对；长度不一时以较短者为准。
func Zip[T, U any](s []T, us []U) []Pair[T, U] {
	n := min(len(s), len(us))
	r := make([]Pair[T, U], n)
	for i := 0; i < n; i++ {
		r[i] = Pair[T, U]{Key: s[i], Val: us[i]}
	}
	return r
}

// Unzip 是 Zip 的逆操作。
func Unzip[T, U any](ps []Pair[T, U]) ([]T, []U) {
	ts := make([]T, len(ps))
	us := make([]U, len(ps))
	for i, p := range ps {
		ts[i], us[i] = p.Key, p.Val
	}
	return ts, us
}
