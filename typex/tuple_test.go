package typex

import (
	"reflect"
	"testing"
)

func TestPair(t *testing.T) {
	p := P("k", 1)
	if p.Key != "k" || p.Val != 1 {
		t.Errorf("P() = {%v %v}, want {k 1}", p.Key, p.Val)
	}
}

func TestZipUnzip(t *testing.T) {
	// 长度不一时以较短者为准
	got := Zip([]string{"a", "b", "c"}, []int{1, 2})
	want := []Pair[string, int]{{Key: "a", Val: 1}, {Key: "b", Val: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Zip() = %v, want %v", got, want)
	}
	ks, vs := Unzip(want)
	if !reflect.DeepEqual(ks, []string{"a", "b"}) || !reflect.DeepEqual(vs, []int{1, 2}) {
		t.Errorf("Unzip() = (%v, %v), want ([a b], [1 2])", ks, vs)
	}
}
