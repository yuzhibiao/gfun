package collection

import (
	"reflect"
	"testing"
)

func TestEntries(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	got := Entries(m)
	if len(got) != 2 {
		t.Fatalf("Entries() len = %d, want 2", len(got))
	}
	for _, p := range got {
		if want, ok := m[p.Key]; !ok || want != p.Val {
			t.Errorf("Entries() contains %v->%v, mismatch with source map", p.Key, p.Val)
		}
	}
}

func TestInvert(t *testing.T) {
	got := Invert(map[string]int{"a": 1, "b": 2})
	want := map[int]string{1: "a", 2: "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Invert() = %v, want %v", got, want)
	}
}

func TestMerge(t *testing.T) {
	got := Merge(
		map[string]int{"a": 1, "b": 2},
		map[string]int{"b": 20, "c": 3},
	)
	want := map[string]int{"a": 1, "b": 20, "c": 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Merge() = %v, want %v", got, want)
	}
}

func TestFilterMap(t *testing.T) {
	got := FilterMap(map[string]int{"a": 1, "b": 2, "c": 3}, func(k string, v int) bool {
		return v > 1
	})
	want := map[string]int{"b": 2, "c": 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FilterMap() = %v, want %v", got, want)
	}
}
