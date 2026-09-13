// Package container 提供泛型数据结构：LRU 缓存、优先级队列、环形缓冲。
// 所有类型均非并发安全，需并发时由调用方自行加锁。
package container

import "container/list"

// LRU 是固定容量的 LRU（最近最少使用）缓存。
type LRU[K comparable, V any] struct {
	cap int
	ll  *list.List // front = 最近使用
	m   map[K]*list.Element
}

type lruEntry[K comparable, V any] struct {
	key K
	val V
}

// NewLRU 创建容量为 capacity 的缓存，capacity 必须 > 0。
func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	if capacity <= 0 {
		panic("container: LRU capacity must be > 0")
	}
	return &LRU[K, V]{
		cap: capacity,
		ll:  list.New(),
		m:   make(map[K]*list.Element, capacity),
	}
}

// Get 取值并把该键移到队头；未命中返回零值和 false。
func (c *LRU[K, V]) Get(k K) (V, bool) {
	if e, ok := c.m[k]; ok {
		c.ll.MoveToFront(e)
		return e.Value.(*lruEntry[K, V]).val, true
	}
	var zero V
	return zero, false
}

// Put 写入；容量满时淘汰队尾（最久未使用）。
func (c *LRU[K, V]) Put(k K, v V) {
	if e, ok := c.m[k]; ok {
		c.ll.MoveToFront(e)
		e.Value.(*lruEntry[K, V]).val = v
		return
	}
	c.m[k] = c.ll.PushFront(&lruEntry[K, V]{key: k, val: v})
	if len(c.m) > c.cap {
		if e := c.ll.Back(); e != nil {
			c.ll.Remove(e)
			delete(c.m, e.Value.(*lruEntry[K, V]).key)
		}
	}
}

// Len 返回当前键值对个数。
func (c *LRU[K, V]) Len() int {
	return len(c.m)
}

// Remove 删除指定键。
func (c *LRU[K, V]) Remove(k K) {
	if e, ok := c.m[k]; ok {
		c.ll.Remove(e)
		delete(c.m, k)
	}
}
