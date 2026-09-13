package container

import (
	"hash/fnv"
	"math"
	"sync"
)

// BloomFilter 是泛型布隆过滤器：用很小的内存判断一个元素是否"可能存在"。
//
// 它只有两种结论：
//   - "一定不存在"：可以放心拒绝；
//   - "可能存在"：需要进一步查实际数据源。
//
// 仅存在假阳性（把不存在的判成可能存在），不存在假阴性（不会漏放已加入的元素）。
// 标准实现不支持删除元素。方法级并发安全。
type BloomFilter[T any] struct {
	mu   sync.RWMutex
	bits []uint64 // 位图，每个 uint64 存 64 位
	m    uint64   // 位总数
	k    int      // 哈希函数个数
	hash func(T) uint64
}

// NewBloomFilter 按"预期元素数 n 和期望误判率 fpRate"创建一个布隆过滤器。
// hash 把 T 转成 64 位哈希；内部用 Kirsch-Mitzenmacher 技巧从这一个哈希
// 派生出 k 个独立的位索引，避免要求调用方提供多个哈希函数。
// n<=0、fpRate 不在 (0,1)、hash 为 nil 时 panic。
func NewBloomFilter[T any](n int, fpRate float64, hash func(T) uint64) *BloomFilter[T] {
	if n <= 0 {
		panic("container: bloom filter n must be > 0")
	}
	if fpRate <= 0 || fpRate >= 1 {
		panic("container: bloom filter fpRate must be in (0, 1)")
	}
	if hash == nil {
		panic("container: bloom filter hash must not be nil")
	}
	m := optimalM(n, fpRate)
	k := optimalK(m, n)
	return &BloomFilter[T]{
		bits: make([]uint64, (m+63)/64),
		m:    m,
		k:    k,
		hash: hash,
	}
}

func optimalM(n int, fp float64) uint64 {
	ln2 := math.Log(2)
	return uint64(math.Ceil(-1 * float64(n) * math.Log(fp) / (ln2 * ln2)))
}

func optimalK(m uint64, n int) int {
	if n <= 0 {
		return 1
	}
	k := int(math.Ceil(float64(m) / float64(n) * math.Log(2)))
	if k < 1 {
		return 1
	}
	return k
}

// Add 把元素加入过滤器。
func (f *BloomFilter[T]) Add(v T) {
	h1, h2 := f.derive(f.hash(v))
	f.mu.Lock()
	for i := 0; i < f.k; i++ {
		idx := (h1 + uint64(i)*h2) % f.m
		f.bits[idx/64] |= 1 << (idx % 64)
	}
	f.mu.Unlock()
}

// MayExist 判断元素是否"可能存在"。返回 false 表示一定不存在，返回 true 表示可能存在。
func (f *BloomFilter[T]) MayExist(v T) bool {
	h1, h2 := f.derive(f.hash(v))
	f.mu.RLock()
	for i := 0; i < f.k; i++ {
		idx := (h1 + uint64(i)*h2) % f.m
		if f.bits[idx/64]&(1<<(idx%64)) == 0 {
			f.mu.RUnlock()
			return false
		}
	}
	f.mu.RUnlock()
	return true
}

// derive 把单一哈希拆成两段，用于 Kirsch-Mitzenmacher 双哈希派生。
func (f *BloomFilter[T]) derive(h uint64) (uint64, uint64) {
	return h, h>>32 | h<<32 // 让两段尽量独立
}

// Reset 清空所有位。
func (f *BloomFilter[T]) Reset() {
	f.mu.Lock()
	for i := range f.bits {
		f.bits[i] = 0
	}
	f.mu.Unlock()
}

// Bits 返回位图总位数；Hashes 返回哈希函数个数。供诊断与持久化使用。
func (f *BloomFilter[T]) Bits() uint64 { return f.m }
func (f *BloomFilter[T]) Hashes() int  { return f.k }

// HashString 用 FNV-1a 返回字符串的 64 位哈希，方便作为 NewBloomFilter 的 hash 参数。
func HashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// HashBytes 用 FNV-1a 返回字节切片的 64 位哈希。
func HashBytes(b []byte) uint64 {
	h := fnv.New64a()
	h.Write(b)
	return h.Sum64()
}
