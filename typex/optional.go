// Package typex 提供类型与值工具：Option（可选值）、Pair（元组）、数值操作。
package typex

// Option 表示"可能有值"，用于区分零值与缺失。
type Option[T any] struct {
	val T
	ok  bool
}

// Some 包装一个存在的值。
func Some[T any](v T) Option[T] {
	return Option[T]{val: v, ok: true}
}

// None 表示值缺失。
func None[T any]() Option[T] {
	return Option[T]{}
}

// IsPresent 报告值是否存在。
func (o Option[T]) IsPresent() bool {
	return o.ok
}

// Get 返回内部值；值不存在时返回 T 的零值，需配合 IsPresent 使用。
func (o Option[T]) Get() T {
	return o.val
}

// OrElse 值存在时返回值，否则返回 fallback。
func (o Option[T]) OrElse(fallback T) T {
	if o.ok {
		return o.val
	}
	return fallback
}

// Unpack 以 Go 惯用的 (value, ok) 形式取出。
func (o Option[T]) Unpack() (T, bool) {
	return o.val, o.ok
}
