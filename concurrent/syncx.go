package concurrent

import (
	"sync"

	"github.com/yuzhibiao/gfun/typex"
)

// Map 是带读写锁的泛型 map，可作 sync.Map 的类型安全替代。
type Map[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewMap 创建空 map。
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{m: make(map[K]V)}
}

// Get 读取键值；未命中时返回零值和 false。
func (m *Map[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.m[k]
	return v, ok
}

// Set 写入或覆盖键值。
func (m *Map[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[k] = v
}

// Delete 删除键。
func (m *Map[K, V]) Delete(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, k)
}

// Len 返回键值对个数。
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.m)
}

// Range 遍历，fn 返回 false 时停止。
// 遍历基于快照进行，fn 内调用本 map 的 Get/Set/Delete 不会死锁。
func (m *Map[K, V]) Range(fn func(K, V) bool) {
	m.mu.RLock()
	pairs := make([]typex.Pair[K, V], 0, len(m.m))
	for k, v := range m.m {
		pairs = append(pairs, typex.Pair[K, V]{Key: k, Val: v})
	}
	m.mu.RUnlock()
	for _, p := range pairs {
		if !fn(p.Key, p.Val) {
			return
		}
	}
}

// Pool 是 sync.Pool 的泛型包装。
type Pool[T any] struct {
	p sync.Pool
}

// NewPool 创建对象池，newFn 用于池空时构造新对象。
func NewPool[T any](newFn func() T) *Pool[T] {
	return &Pool[T]{p: sync.Pool{New: func() any { return newFn() }}}
}

// Get 取出一个对象，池空时调用 newFn 构造。
func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}

// Put 归还对象。
func (p *Pool[T]) Put(v T) {
	p.p.Put(v)
}
