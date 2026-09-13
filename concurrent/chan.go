// Package concurrent 提供并发辅助：channel 工具、泛型并发 map、对象池、重试。
package concurrent

// Generate 把一组值发送到带缓冲的 channel 后立即关闭。
// 缓冲区大小等于元素个数，因此不存在后台 goroutine，
// 消费方提前退出也不会有任何泄漏。
func Generate[T any](vals ...T) <-chan T {
	out := make(chan T, len(vals))
	for _, v := range vals {
		out <- v // 缓冲足够，永不阻塞
	}
	close(out)
	return out
}

// Take 从 ch 接收最多 n 个值并收集为切片，ch 关闭时提前结束。
func Take[T any](ch <-chan T, n int) []T {
	r := make([]T, 0, n)
	for i := 0; i < n; i++ {
		v, ok := <-ch
		if !ok {
			break
		}
		r = append(r, v)
	}
	return r
}
