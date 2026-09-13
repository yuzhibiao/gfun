package concurrent

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestGenerateCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := GenerateCtx(ctx, 1, 2, 3)
	if got := Take(ch, 1); !reflect.DeepEqual(got, []int{1}) {
		t.Errorf("Take() = %v, want [1]", got)
	}
	cancel()
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("channel should be closed after cancel")
		}
	case <-time.After(time.Second):
		t.Error("channel not closed within 1s after cancel")
	}
}

func TestGenerateTake(t *testing.T) {
	ch := Generate(1, 2, 3)
	if got := Take(ch, 2); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Errorf("Take() = %v, want [1 2]", got)
	}
	if got := Take(ch, 5); !reflect.DeepEqual(got, []int{3}) {
		t.Errorf("Take() after close = %v, want [3]", got)
	}
}

func TestGenerateEmpty(t *testing.T) {
	ch := Generate[int]()
	select {
	case v, ok := <-ch:
		if ok {
			t.Errorf("expected closed channel, got %v", v)
		}
	case <-time.After(time.Second):
		t.Error("channel not closed within 1s")
	}
}
