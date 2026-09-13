package concurrent

import (
	"errors"
	"testing"
	"time"
)

func TestRetrySuccess(t *testing.T) {
	calls := 0
	err := Retry(5, 0, func() error {
		calls++
		if calls < 3 {
			return errors.New("fail")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Retry() err = %v, want nil", err)
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}

func TestRetryExhausted(t *testing.T) {
	calls := 0
	err := Retry(3, 0, func() error {
		calls++
		return errors.New("always fail")
	})
	if err == nil {
		t.Fatal("Retry() err = nil, want error")
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}

func TestRetryWithBackoff(t *testing.T) {
	calls := 0
	start := time.Now()
	_ = RetryWithBackoff(3, time.Millisecond, func() error {
		calls++
		return errors.New("always fail")
	})
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
	// 两次退避：1ms + 2ms = 3ms 起步
	if elapsed := time.Since(start); elapsed < 2*time.Millisecond {
		t.Errorf("elapsed = %v, want >= 3ms", elapsed)
	}
}
