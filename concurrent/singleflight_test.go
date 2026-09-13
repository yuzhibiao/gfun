package concurrent

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingleFlightDedup(t *testing.T) {
	g := NewGroup[string, int]()
	var calls int32
	const n = 1000
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			v, err := g.Do("k", func() (int, error) {
				atomic.AddInt32(&calls, 1)
				time.Sleep(10 * time.Millisecond) // 让 1000 个并发真实重叠
				return 42, nil
			})
			if err != nil || v != 42 {
				t.Errorf("got (%d, %v), want (42, nil)", v, err)
			}
		}()
	}
	close(start)
	wg.Wait()
	if calls := atomic.LoadInt32(&calls); calls != 1 {
		t.Errorf("fn 执行了 %d 次，期望 1 次（应去重）", calls)
	}
}

func TestSingleFlightDifferentKeys(t *testing.T) {
	g := NewGroup[int, int]()
	var calls int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, _ = g.Do(i, func() (int, error) {
				atomic.AddInt32(&calls, 1)
				return i, nil
			})
		}(i)
	}
	wg.Wait()
	if calls := atomic.LoadInt32(&calls); calls != 100 {
		t.Errorf("不同 key 应各自执行，期望 100 次，实际 %d 次", calls)
	}
}

func TestSingleFlightErrorPropagation(t *testing.T) {
	g := NewGroup[string, int]()
	sentinel := errors.New("boom")
	var wg sync.WaitGroup
	const n = 50
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, errs[i] = g.Do("k", func() (int, error) {
				return 0, sentinel
			})
		}(i)
	}
	wg.Wait()
	for i, e := range errs {
		if !errors.Is(e, sentinel) {
			t.Fatalf("第 %d 个调用未拿到共享错误，got %v", i, e)
		}
	}
}

func TestSingleFlightConcurrentMixed(t *testing.T) {
	g := NewGroup[string, int]()
	var calls int32
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('a' + i%10))
			_, _ = g.Do(key, func() (int, error) {
				atomic.AddInt32(&calls, 1)
				time.Sleep(5 * time.Millisecond) // 让同 key 的并发真实重叠
				return i, nil
			})
		}(i)
	}
	wg.Wait()
	// 10 个 key、1000 个调用，去重后应远少于 1000；-race 下调度有抖动，留足余量
	if c := atomic.LoadInt32(&calls); c >= 200 {
		t.Errorf("去重后应远少于 1000 次，实际 %d 次", c)
	}
}
