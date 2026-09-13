package collection

import "sync"

// Set 是基于 map 实现的集合。
// 所有方法并发安全（单方法原子）；复合操作（如先 Has 再 Add）
// 不保证原子性，需要时由调用方自行加锁。
type Set[T comparable] struct {
	mu sync.RWMutex
	m  map[T]struct{}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range vals {
		s.m[v] = struct{}{}
	}
}

// Has 判断元素是否存在。
func (s *Set[T]) Has(v T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[v]
	return ok
}

// Delete 删除元素，元素不存在时无操作。
func (s *Set[T]) Delete(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, v)
}

// Len 返回元素个数。
func (s *Set[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}

// ToSlice 返回所有元素的快照，顺序不确定。
func (s *Set[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r := make([]T, 0, len(s.m))
	for v := range s.m {
		r = append(r, v)
	}
	return r
}

// Union 返回 s 与 others 的并集。
// 通过 ToSlice 快照读取各方元素，不存在嵌套加锁，并发安全。
func (s *Set[T]) Union(others ...*Set[T]) *Set[T] {
	vals := s.ToSlice()
	for _, o := range others {
		vals = append(vals, o.ToSlice()...)
	}
	r := NewSet[T]()
	r.Add(vals...)
	return r
}

// Intersect 返回 s 与 other 的交集。
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	r := NewSet[T]()
	for _, v := range s.ToSlice() { // 快照，避免持锁调用 other
		if other.Has(v) {
			r.Add(v)
		}
	}
	return r
}

// Difference 返回在 s 中但不在 other 中的元素。
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	r := NewSet[T]()
	for _, v := range s.ToSlice() { // 快照，避免持锁调用 other
		if !other.Has(v) {
			r.Add(v)
		}
	}
	return r
}
