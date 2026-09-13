// Package concurrent 提供并发辅助：channel 工具、泛型并发 map、对象池、重试。
package concurrent

import "context"

// Generate 把一组值依次发送到返回的 channel，发送完毕后关闭。
// 注意：消费方提前退出时（未读完即不再接收），生产 goroutine 会阻塞在
// 发送上造成泄漏，此类场景请使用 GenerateCtx。
func Generate[T any](vals ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, v := range vals {
			out <- v
		}
	}()
	return out
}

// GenerateCtx 是带取消的 Generate：ctx 取消后停止发送并关闭 channel，
// 用于消费方可能提前退出的场景，避免 goroutine 泄漏。
func GenerateCtx[T any](ctx context.Context, vals ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, v := range vals {
			select {
			case out <- v:
			case <-ctx.Done():
				return
			}
		}
	}()
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
