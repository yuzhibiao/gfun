package container

import "container/heap"

// PriorityQueue 是基于二叉堆的优先级队列。
type PriorityQueue[T any] struct {
	h *pqHeap[T]
}

type pqHeap[T any] struct {
	s    []T
	less func(a, b T) bool
}

func (h *pqHeap[T]) Len() int           { return len(h.s) }
func (h *pqHeap[T]) Less(i, j int) bool { return h.less(h.s[i], h.s[j]) }
func (h *pqHeap[T]) Swap(i, j int)      { h.s[i], h.s[j] = h.s[j], h.s[i] }
func (h *pqHeap[T]) Push(x any)         { h.s = append(h.s, x.(T)) }
func (h *pqHeap[T]) Pop() any {
	n := len(h.s)
	v := h.s[n-1]
	h.s = h.s[:n-1]
	return v
}

// NewPriorityQueue 创建空队列；less(a, b) 为 true 表示 a 先于 b 出队。
func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{h: &pqHeap[T]{less: less}}
}

// Push 压入一个元素。
func (q *PriorityQueue[T]) Push(v T) {
	heap.Push(q.h, v)
}

// Pop 弹出最高优先级元素；队列为空时 panic（与 container/heap 行为一致）。
func (q *PriorityQueue[T]) Pop() T {
	return heap.Pop(q.h).(T)
}

// Peek 查看队首元素但不弹出；队列为空时 panic。
func (q *PriorityQueue[T]) Peek() T {
	return q.h.s[0]
}

// Len 返回元素个数。
func (q *PriorityQueue[T]) Len() int {
	return q.h.Len()
}
