package concurrent

import (
	"reflect"
	"testing"
	"time"
)

func TestGenerateTake(t *testing.T) {
	ch := Generate(1, 2, 3)
	if got := Take(ch, 2); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Errorf("Take() = %v, want [1 2]", got)
	}
	// 消费方提前退出：剩余值留在缓冲里，无 goroutine 泄漏
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
