package container

import (
	"reflect"
	"sync"
	"testing"
)

func TestRingConcurrent(t *testing.T) {
	r := NewRing[int](64)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				r.Push(g*100 + i)
				r.Len()
				if i%10 == 0 {
					r.All()
				}
			}
		}(g)
	}
	wg.Wait()
	if r.Len() != 64 {
		t.Errorf("Len() = %d, want 64", r.Len())
	}
}

func TestRingUnderfill(t *testing.T) {
	r := NewRing[int](5)
	r.Push(1)
	r.Push(2)
	if r.Len() != 2 {
		t.Errorf("Len() = %d, want 2", r.Len())
	}
	if got := r.All(); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Errorf("All() = %v, want [1 2]", got)
	}
}

func TestRingOverwrite(t *testing.T) {
	r := NewRing[int](3)
	for _, v := range []int{1, 2, 3, 4, 5} {
		r.Push(v)
	}
	if r.Len() != 3 {
		t.Errorf("Len() = %d, want 3", r.Len())
	}
	// 写满后 4、5 覆盖了 1、2，按旧到新顺序
	if got := r.All(); !reflect.DeepEqual(got, []int{3, 4, 5}) {
		t.Errorf("All() = %v, want [3 4 5]", got)
	}
}

func TestRingCapacityOne(t *testing.T) {
	r := NewRing[string](1)
	r.Push("a")
	r.Push("b")
	if got := r.All(); !reflect.DeepEqual(got, []string{"b"}) {
		t.Errorf("All() = %v, want [b]", got)
	}
}
