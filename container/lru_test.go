package container

import "testing"

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
