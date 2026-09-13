package collection

// Set 是基于 map 实现的集合，非并发安全。
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet 用初始元素创建集合。
func NewSet[T comparable](vals ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(vals))}
	for _, v := range vals {
		s.m[v] = struct{}{}
	}
	return s
}

// Add 添加元素。
func (s *Set[T]) Add(vals ...T) {
	for _, v := range vals {
		s.m[v] = struct{}{}
	}
}

// Has 判断元素是否存在。
func (s *Set[T]) Has(v T) bool {
	_, ok := s.m[v]
	return ok
}

// Delete 删除元素，元素不存在时无操作。
func (s *Set[T]) Delete(v T) {
	delete(s.m, v)
}

// Len 返回元素个数。
func (s *Set[T]) Len() int {
	return len(s.m)
}

// ToSlice 返回所有元素，顺序不确定。
func (s *Set[T]) ToSlice() []T {
	r := make([]T, 0, len(s.m))
	for v := range s.m {
		r = append(r, v)
	}
	return r
}

// Union 返回 s 与 others 的并集。
func (s *Set[T]) Union(others ...*Set[T]) *Set[T] {
	r := NewSet[T]()
	r.Add(s.ToSlice()...)
	for _, o := range others {
		r.Add(o.ToSlice()...)
	}
	return r
}

// Intersect 返回 s 与 other 的交集。
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	r := NewSet[T]()
	for v := range s.m {
		if other.Has(v) {
			r.Add(v)
		}
	}
	return r
}

// Difference 返回在 s 中但不在 other 中的元素。
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	r := NewSet[T]()
	for v := range s.m {
		if !other.Has(v) {
			r.Add(v)
		}
	}
	return r
}
