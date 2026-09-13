package container

import "testing"

func TestPriorityQueueInt(t *testing.T) {
	q := NewPriorityQueue(func(a, b int) bool { return a < b }) // 小顶堆
	q.Push(5)
	q.Push(1)
	q.Push(3)
	if q.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", q.Len())
	}
	if q.Peek() != 1 {
		t.Errorf("Peek() = %d, want 1", q.Peek())
	}
	var got []int
	for q.Len() > 0 {
		got = append(got, q.Pop())
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 3 || got[2] != 5 {
		t.Errorf("pop order = %v, want [1 3 5]", got)
	}
}

func TestPriorityQueueCustom(t *testing.T) {
	type task struct {
		name     string
		priority int
	}
	q := NewPriorityQueue(func(a, b task) bool { return a.priority > b.priority }) // 大顶
	q.Push(task{"low", 1})
	q.Push(task{"high", 9})
	if first := q.Pop(); first.name != "high" {
		t.Errorf("Pop() = %s, want high", first.name)
	}
}
