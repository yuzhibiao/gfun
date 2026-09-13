package concurrent

import "sync"

type call[V any] struct {
	wg  sync.WaitGroup
	val V
	err error
}

// Group 是泛型版 singleflight：对同一 key 的并发 Do 调用只执行一次 fn，
// 其他调用阻塞等待并共享结果。用于消除"同一份热数据的并发重建"。
//
// 注意：fn 不应 panic；若 fn panic，等待方会读到零值与 nil 错误，
// 且 panic 会沿调用栈向上抛出。
type Group[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]*call[V]
}

// NewGroup 创建空的 singleflight 组。零值 Group 也可直接使用（Do 内部懒初始化 m）。
func NewGroup[K comparable, V any]() *Group[K, V] {
	return &Group[K, V]{m: make(map[K]*call[V])}
}

// Do 对 key k 执行 fn。若已有同 key 的调用在进行中，则阻塞等待其结果而不重复执行 fn。
func (g *Group[K, V]) Do(k K, fn func() (V, error)) (V, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[K]*call[V])
	}
	if c, ok := g.m[k]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := &call[V]{}
	c.wg.Add(1)
	g.m[k] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, k)
	g.mu.Unlock()

	return c.val, c.err
}
