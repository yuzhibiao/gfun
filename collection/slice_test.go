package collection

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(v int) int { return v * 2 })
	if want := []int{2, 4, 6}; !reflect.DeepEqual(got, want) {
		t.Errorf("Map() = %v, want %v", got, want)
	}
	if got := Map(nil, func(v int) int { return v }); len(got) != 0 {
		t.Errorf("Map(nil) = %v, want empty", got)
	}
}

func TestFilter(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	if want := []int{2, 4}; !reflect.DeepEqual(got, want) {
		t.Errorf("Filter() = %v, want %v", got, want)
	}
}

func TestReduce(t *testing.T) {
	if got := Reduce([]int{1, 2, 3, 4}, 0, func(acc, v int) int { return acc + v }); got != 10 {
		t.Errorf("Reduce() = %d, want 10", got)
	}
	if got := Reduce([]string{"a", "b"}, "", func(acc, v string) string { return acc + v }); got != "ab" {
		t.Errorf("Reduce() = %q, want ab", got)
	}
}

func TestChunk(t *testing.T) {
	got := Chunk([]int{1, 2, 3, 4, 5}, 2)
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Chunk() = %v, want %v", got, want)
	}
	if got := Chunk([]int{1, 2}, 0); got != nil {
		t.Errorf("Chunk(size=0) = %v, want nil", got)
	}
}

func TestUnique(t *testing.T) {
	got := Unique([]string{"a", "b", "a", "c", "b"})
	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Unique() = %v, want %v", got, want)
	}
}

func TestGroupBy(t *testing.T) {
	got := GroupBy([]int{1, 2, 3, 4}, func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	})
	want := map[string][]int{"odd": {1, 3}, "even": {2, 4}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupBy() = %v, want %v", got, want)
	}
}
