package concurrent

import "time"

// Retry 重试 fn 直到成功或达到 attempts 次（首次调用计入 attempts）。
// 每次失败后间隔 interval 再重试，最后一次失败后立即返回错误。
// attempts <= 0 时视为 1。
func Retry(attempts int, interval time.Duration, fn func() error) error {
	if attempts <= 0 {
		attempts = 1
	}
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < attempts-1 && interval > 0 {
			time.Sleep(interval)
		}
	}
	return err
}

// RetryWithBackoff 指数退避重试：第 n 次失败（n 从 1 起）后等待 base * 2^(n-1)。
func RetryWithBackoff(attempts int, base time.Duration, fn func() error) error {
	if attempts <= 0 {
		attempts = 1
	}
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < attempts-1 && base > 0 {
			time.Sleep(base << uint(i))
		}
	}
	return err
}
