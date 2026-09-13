package concurrent

import (
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Errorf("Get() = (%d, %v), want (1, true)", v, ok)
	}
	if _, ok := m.Get("missing"); ok {
		t.Error("Get(missing) ok = true, want false")
	}
	if m.Len() != 1 {
		t.Errorf("Len() = %d, want 1", m.Len())
	}
	m.Delete("a")
	if m.Len() != 0 {
		t.Errorf("Len() after Delete = %d, want 0", m.Len())
	}
}

func TestMapConcurrent(t *testing.T) {
	m := NewMap[int, int]()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			m.Set(n, n*2)
		}(i)
	}
	wg.Wait()
	if m.Len() != 50 {
		t.Errorf("Len() = %d, want 50", m.Len())
	}
}

func TestMapRangeModify(t *testing.T) {
	m := NewMap[int, int]()
	for i := 0; i < 5; i++ {
		m.Set(i, i)
	}
	// 旧实现在持 RLock 状态下回调 Set 会死锁；快照版应能正常工作
	m.Range(func(k, v int) bool {
		m.Set(k, v*10)
		return true
	})
	for i := 0; i < 5; i++ {
		if v, _ := m.Get(i); v != i*10 {
			t.Errorf("Get(%d) = %d, want %d", i, v, i*10)
		}
	}
}

func TestPool(t *testing.T) {
	p := NewPool(func() []byte { return make([]byte, 8) })
	buf := p.Get()
	if len(buf) != 8 {
		t.Errorf("Get() len = %d, want 8", len(buf))
	}
	p.Put(buf)
	if again := p.Get(); len(again) != 8 {
		t.Errorf("Get() after Put len = %d, want 8", len(again))
	}
}
