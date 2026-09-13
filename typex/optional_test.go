package typex

import "testing"

func TestSome(t *testing.T) {
	o := Some(42)
	if !o.IsPresent() {
		t.Error("Some().IsPresent() = false, want true")
	}
	if o.Get() != 42 {
		t.Errorf("Get() = %d, want 42", o.Get())
	}
	if v, ok := o.Unpack(); !ok || v != 42 {
		t.Errorf("Unpack() = (%d, %v), want (42, true)", v, ok)
	}
}

func TestNone(t *testing.T) {
	o := None[int]()
	if o.IsPresent() {
		t.Error("None().IsPresent() = true, want false")
	}
	if o.Get() != 0 {
		t.Errorf("Get() = %d, want 0 (zero value)", o.Get())
	}
	if got := o.OrElse(7); got != 7 {
		t.Errorf("OrElse(7) = %d, want 7", got)
	}
	if got := Some(1).OrElse(7); got != 1 {
		t.Errorf("Some(1).OrElse(7) = %d, want 1", got)
	}
}
