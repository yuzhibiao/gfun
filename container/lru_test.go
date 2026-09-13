package container

import (
	"sync"
	"testing"
)

func TestLRUConcurrent(t *testing.T) {
	c := NewLRU[int, int](50)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				k := (g*200 + i) % 40
				c.Put(k, i)
				c.Get(k)
				c.Len()
				if i%10 == 0 {
					c.Remove(k)
				}
			}
		}(g)
	}
	wg.Wait()
}

func TestLRUPutGet(t *testing.T) {
	c := NewLRU[string, int](2)
	c.Put("a", 1)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Errorf("Get() = (%d, %v), want (1, true)", v, ok)
	}
	if _, ok := c.Get("missing"); ok {
		t.Error("Get(missing) ok = true, want false")
	}
}

func TestLRUEviction(t *testing.T) {
	c := NewLRU[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a")    // a 变为最近使用
	c.Put("c", 3) // 容量满，淘汰 b
	if _, ok := c.Get("b"); ok {
		t.Error("b should be evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Error("a should survive")
	}
	if c.Len() != 2 {
		t.Errorf("Len() = %d, want 2", c.Len())
	}
}

func TestLRURemove(t *testing.T) {
	c := NewLRU[int, int](2)
	c.Put(1, 1)
	c.Remove(1)
	if c.Len() != 0 {
		t.Errorf("Len() after Remove = %d, want 0", c.Len())
	}
}
