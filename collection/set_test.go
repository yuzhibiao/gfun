package collection

import "testing"

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
