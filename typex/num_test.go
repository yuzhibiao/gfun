package typex

import "testing"

func TestClamp(t *testing.T) {
	tests := []struct {
		v, lo, hi, want int
	}{
		{5, 1, 10, 5},
		{0, 1, 10, 1},
		{11, 1, 10, 10},
	}
	for _, tt := range tests {
		if got := Clamp(tt.v, tt.lo, tt.hi); got != tt.want {
			t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.v, tt.lo, tt.hi, got, tt.want)
		}
	}
	if got := Clamp("m", "a", "f"); got != "f" {
		t.Errorf("Clamp(string) = %q, want %q", got, "f")
	}
}

func TestSum(t *testing.T) {
	if got := Sum(1, 2, 3, 4); got != 10 {
		t.Errorf("Sum() = %d, want 10", got)
	}
	if got := Sum[int](); got != 0 {
		t.Errorf("Sum() = %d, want 0", got)
	}
	if got := Sum("a", "b"); got != "ab" {
		t.Errorf("Sum(string) = %q, want %q", got, "ab")
	}
}
