package container

// Ring 是固定容量的环形缓冲，写满后新数据覆盖最旧数据。
// 适合滑动窗口、限流统计等场景。
type Ring[T any] struct {
	buf  []T
	head int // 下一个写入位置
	full bool
}

// NewRing 创建容量为 capacity 的环形缓冲，capacity 必须 > 0。
func NewRing[T any](capacity int) *Ring[T] {
	if capacity <= 0 {
		panic("container: Ring capacity must be > 0")
	}
	return &Ring[T]{buf: make([]T, capacity)}
}

// Push 写入一个元素。
func (r *Ring[T]) Push(v T) {
	r.buf[r.head] = v
	r.head = (r.head + 1) % len(r.buf)
	if r.head == 0 {
		r.full = true
	}
}

// Len 返回当前元素个数。
func (r *Ring[T]) Len() int {
	if r.full {
		return len(r.buf)
	}
	return r.head
}

// All 按写入顺序（旧到新）返回所有元素的副本。
func (r *Ring[T]) All() []T {
	n := r.Len()
	out := make([]T, 0, n)
	start := 0
	if r.full {
		start = r.head // 写满后最旧元素在 head 处
	}
	for i := 0; i < n; i++ {
		out = append(out, r.buf[(start+i)%len(r.buf)])
	}
	return out
}
