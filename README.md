# gfun

基于 Go 泛型的常用工具库，提供集合操作、类型工具、并发辅助与泛型数据结构，以及数据库列类型转换（JSON 列 / 加密列）。

- 模块路径：`github.com/yuzhibiao/gfun`
- Go 版本要求：**1.21+**
- 标准库 `slices` / `maps` / `cmp` 已覆盖的能力（排序、查找、克隆等）不重复实现，只补缺口
- 所有代码通过 `gofmt` / `go vet` / `go test -race`

## 安装

```bash
go get github.com/yuzhibiao/gfun@v0.1.0
```

## 包结构

| 包 | 内容 |
|---|---|
| `collection` | 切片 / map / 集合的泛型操作 |
| `typex` | 类型与值工具：`Option[T]`、`Pair`、数值操作 |
| `concurrent` | channel 工具、泛型并发 map、对象池、重试 |
| `container` | 泛型数据结构：LRU、优先级队列、环形缓冲 |
| `sqlx` | 数据库列类型转换：`JsonColumn[T]`、`EncryptColumn[T]` |

## 使用示例

### collection

```go
import "github.com/yuzhibiao/gfun/collection"

double := collection.Map([]int{1, 2, 3}, func(v int) int { return v * 2 })
// [2 4 6]

evens := collection.Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
// [2 4]

total := collection.Reduce([]int{1, 2, 3}, 0, func(acc, v int) int { return acc + v })
// 6

parts := collection.Chunk([]int{1, 2, 3, 4, 5}, 2)
// [[1 2] [3 4] [5]]

groups := collection.GroupBy(users, func(u User) string {
    return u.City
}) // map[城市][]User

s := collection.NewSet(1, 2, 3)
s.Intersect(collection.NewSet(2, 3, 4)) // {2, 3}
```

### typex

```go
import "github.com/yuzhibiao/gfun/typex"

// Option：区分零值与缺失
o := typex.Some(42)
v := o.OrElse(0) // 42

// Pair / Zip
pairs := typex.Zip([]string{"a", "b"}, []int{1, 2})
// [{a 1} {b 2}]

// 数值
clamped := typex.Clamp(score, 0, 100)
```

### concurrent

```go
import "github.com/yuzhibiao/gfun/concurrent"

// channel 工具
ch := concurrent.Generate(1, 2, 3)
vals := concurrent.Take(ch, 2) // [1 2]，内部无 goroutine，提前退出不泄漏

// 泛型并发 map（sync.Map 的类型安全替代）
m := concurrent.NewMap[string, int]()
m.Set("a", 1)
v, ok := m.Get("a")

// 重试（固定间隔 / 指数退避）
err := concurrent.Retry(3, time.Second, func() error {
    return callRemoteAPI()
})
```

### container

```go
import "github.com/yuzhibiao/gfun/container"

// LRU 缓存
cache := container.NewLRU[string, any](100)
cache.Put("k", v)
v, ok := cache.Get("k")

// 优先级队列（less 为真者先出队）
pq := container.NewPriorityQueue(func(a, b int) bool { return a < b }) // 小顶堆
pq.Push(3)
pq.Push(1)
pq.Pop() // 1

// 环形缓冲（滑动窗口）
ring := container.NewRing[float64](60) // 最近 60 秒
ring.Push(0.85)
```

### sqlx

```go
import "github.com/yuzhibiao/gfun/sqlx"

type User struct {
    Ext   sqlx.JsonColumn[map[string]any] `gorm:"type:json"`   // JSON 列
    Phone sqlx.EncryptColumn[string]       `gorm:"type:text"`   // 加密列
}

// EncryptColumn 需在初始化阶段设置一次全局密钥（16/24/32 字节）
func init() {
    key := []byte(os.Getenv("GFUN_ENCRYPT_KEY")) // 密钥务必来自环境变量等外部配置，不要硬编码
    if err := sqlx.SetEncryptKey(key); err != nil {
        panic(err)
    }
}
```

- `JsonColumn`：写库时自动 `json.Marshal`，读库时自动 `Unmarshal`，`Valid=false` 时写 `NULL`
- `EncryptColumn`：写库时 AES-GCM 加密后 base64 存储，读库时解密还原；密文被篡改会解密失败（GCM 认证）

## 运行测试

```bash
go test -race ./...
```

## 注意事项

- `collection` 的 `Map` / `Filter` / `Chunk` 等均为值语义，不修改入参；`Chunk` 返回的子切片与原切片共享底层数组
- 所有类型（`Set`、`concurrent.Map`、`LRU`、`PriorityQueue`、`Ring` 等）均为**方法级并发安全**：单个方法可并发调用，但"先查再改"这类复合操作不保证原子性，需要时由调用方自行加锁
- `concurrent.Generate` 内部无后台 goroutine（缓冲等于元素个数），消费方提前退出无泄漏
- `sqlx.EncryptColumn` 的密钥管理遵循"密钥永不进代码库"原则
