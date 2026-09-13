package container

import (
	"fmt"
	"sync"
	"testing"
)

func TestBloomFilterNoFalseNegative(t *testing.T) {
	f := NewBloomFilter[string](10000, 0.01, HashString)
	for i := 0; i < 10000; i++ {
		f.Add(fmt.Sprintf("item-%d", i))
	}
	for i := 0; i < 10000; i++ {
		if !f.MayExist(fmt.Sprintf("item-%d", i)) {
			t.Fatalf("item-%d 应已加入却被判为不存在（假阴性）", i)
		}
	}
}

func TestBloomFilterFalsePositiveRate(t *testing.T) {
	const n = 10000
	fp := 0.01
	f := NewBloomFilter[int](n, fp, func(i int) uint64 { return uint64(i)*2654435761 + 1 })

	for i := 0; i < n; i++ {
		f.Add(i)
	}
	// 用 [n, n+10000) 这批"没加入"的元素统计假阳性率
	var fpCount int
	for i := n; i < 2*n; i++ {
		if f.MayExist(i) {
			fpCount++
		}
	}
	rate := float64(fpCount) / float64(n)
	// 理论 ~1%，给 4 倍余量避免概率性测试 flaky
	if rate > 4*fp {
		t.Errorf("假阳性率 %.3f 超出 %.3f", rate, 4*fp)
	}
}

func TestBloomFilterReset(t *testing.T) {
	f := NewBloomFilter[string](100, 0.01, HashString)
	f.Add("a")
	f.Add("b")
	if !f.MayExist("a") || !f.MayExist("b") {
		t.Fatal("加入后应能查到")
	}
	f.Reset()
	if f.MayExist("a") || f.MayExist("b") {
		t.Fatal("Reset 后应判为不存在")
	}
}

func TestBloomFilterPanicsOnBadParams(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("n<=0 应 panic")
		}
	}()
	_ = NewBloomFilter[string](0, 0.01, HashString)
}

func TestBloomFilterConcurrent(t *testing.T) {
	f := NewBloomFilter[int](10000, 0.01, func(i int) uint64 { return uint64(i) })
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				v := g*1000 + i
				f.Add(v)
				_ = f.MayExist(v)
			}
		}(g)
	}
	wg.Wait()
	// 全部加入后应都能查到
	for g := 0; g < 8; g++ {
		for i := 0; i < 1000; i++ {
			if !f.MayExist(g*1000 + i) {
				t.Fatal("并发加入后存在假阴性")
			}
		}
	}
}
