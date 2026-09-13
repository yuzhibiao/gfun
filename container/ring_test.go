package container

import (
	"reflect"
	"testing"
)

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
