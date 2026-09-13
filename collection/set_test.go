package collection

import (
	"sync"
	"testing"
)

func TestSetConcurrent(t *testing.T) {
	s := NewSet[int]()
	a := NewSet(10, 11)
	b := NewSet(11, 12)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				k := (g*100 + i) % 30
				s.Add(k)
				s.Has(k)
				s.Len()
				s.ToSlice()
				s.Intersect(a)
				s.Union(a, b) // 并发读 a/b 并发写 s，验证无嵌套锁死锁
				if i%5 == 0 {
					s.Delete(k)
				}
			}
		}(g)
	}
	wg.Wait()
}

func TestSetBasic(t *testing.T) {
	s := NewSet(1, 2)
	s.Add(3, 3) // 重复添加
	if s.Len() != 3 {
		t.Errorf("Len() = %d, want 3", s.Len())
	}
	if !s.Has(2) {
		t.Error("Has(2) = false, want true")
	}
	s.Delete(2)
	if s.Has(2) {
		t.Error("Has(2) after Delete = true, want false")
	}
}

func TestSetOps(t *testing.T) {
	a := NewSet(1, 2, 3)
	b := NewSet(2, 3, 4)

	if got := a.Intersect(b).ToSlice(); !sameElements(got, []int{2, 3}) {
		t.Errorf("Intersect() = %v, want {2,3}", got)
	}
	if got := a.Difference(b).ToSlice(); !sameElements(got, []int{1}) {
		t.Errorf("Difference() = %v, want {1}", got)
	}
	if got := a.Union(b).ToSlice(); !sameElements(got, []int{1, 2, 3, 4}) {
		t.Errorf("Union() = %v, want {1,2,3,4}", got)
	}
}

// sameElements 忽略顺序比较两个切片的元素集合。
func sameElements[T comparable](got, want []T) bool {
	if len(got) != len(want) {
		return false
	}
	m := make(map[T]struct{}, len(got))
	for _, v := range got {
		m[v] = struct{}{}
	}
	for _, v := range want {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}
