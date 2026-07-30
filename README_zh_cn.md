# 🐱 rainbowcat

Go 工具库。

[🇬🇧 English](README.md)

## 安装

```bash
go get github.com/rambollwong/rainbowcat
```

需要 **Go 1.26+**。

## 包概览

| 包 | 说明 |
|---|---|
| [`cache`](#cache) | 泛型 FIFO 缓存，支持可选的线程安全和淘汰回调 |
| [`global`](#global) | 全局 `sync.WaitGroup`，用于优雅关闭协调 |
| [`logo`](#logo) | ASCII 艺术横幅打印（装饰性） |
| [`pipeline`](#pipeline) | 多阶段并行任务流水线 — 每阶段并发执行 |
| [`pool`](#pool) | Goroutine 工作池 + `[]byte` 复用池 |
| [`signal`](#signal) | 操作系统信号处理（`SIGINT`/`SIGTERM`），含平台变体 |
| [`smtp`](#smtp) | SMTP 邮件发送器（封装 `gomail`） |
| [`task`](#task) | 基于 timer/ticker 的任务监控器，支持可插拔 `DataStore` |
| [`time`](#time) | 单一辅助函数：`UseLocalUTC()` |
| [`types`](#types) | 泛型辅助：`Set[T]`、`Entry[K,V]`、`Clonable[T]` |
| [`util`](#util) | 切片与 map 操作、字节/十六进制转换、gzip、斐波那契、字符串解析 |
| [`writer/filewriter`](#writerfilewriter) | 滚动文件写入器：按大小和按时间 |

---

### cache

基于 `container/list` + `map` 的泛型固定容量 FIFO 缓存，支持可选的 `sync.RWMutex` 线程安全和淘汰回调。

```go
c := cache.NewFIFOCache[int, string](100, true)
c.Put(1, "hello")
v, ok := c.Get(1) // "hello", true
c.Size()           // 1
```

### global

进程级 `sync.WaitGroup`，用于优雅关闭时跟踪运行中的 goroutine。

```go
global.RunTask(func() {
    // 执行任务
})
// ... 等待所有任务完成
global.Wait()
```

### pool

**Worker pool** — 有界 goroutine 池，支持函数式选项：

```go
wp := pool.NewWorkerPool(4,
    pool.WithBufferSize(16),
    pool.WithRejectHandler(func(t pool.Task) { /* 日志记录 */ }),
)
err := wp.Submit(func() { /* 执行任务 */ })
wp.Close()
```

**Bytes pool** — 封装 `sync.Pool`，用于 `[]byte` 复用：

```go
bp := pool.NewBytesPool(512, pool.DefaultMaxBytesCap)
bz := bp.Get()
defer bp.Put(bz)
```

### pipeline

多阶段并发流水线。每阶段可配置并发数；第 N 阶段的输出成为第 N+1 阶段的输入。

```go
stage1 := pipeline.GenericTaskProvider[int, string](func(i int) (string, bool) {
    return fmt.Sprintf("got %d", i), true
})
pp, _ := pipeline.RunParallelTaskPipeline(2, []uint8{2, 1}, stage1)
pp.PushJob(42)
fmt.Println(<-pp.OutputC()) // "got 42"
pp.Close()
```

### signal

监听 `SIGINT`/`SIGTERM` 并执行回调：

```go
signal.WatchExitSignal(func() {
    fmt.Println("正在关闭...")
})
signal.WaitForKeyPress() // 跨平台"按任意键继续"
```

### task

基于 timer/ticker 的任务调度，支持可插拔 `DataStore`：

```go
tm := &task.TasksMonitor{}
tm.SetDataStore(myDataStore)
tm.RegisterTickerForTasks(time.Minute, "cleanup", func(d task.Data) {
    // 定期清理
})
tm.Start()
```

### types

线程安全的泛型 `Set[T]`，基于 `map + sync.RWMutex`：

```go
s := types.NewSet[int]()
s.Put(1)      // true
s.Put(1)      // false (已存在)
s.Exist(1)    // true
s.Size()      // 1
s.Range(func(v int) bool { fmt.Println(v); return true })
```

### util

丰富的切片与 map 工具函数（受 `samber/lo` 启发），以及 gzip 压缩、斐波那契数列、字节转换和字符串解析：

```go
// 切片操作
util.SliceFilter([]int{1, 2, 3, 4}, func(i int, v int) bool { return v%2 == 0 })
util.SliceShuffle([]int{1, 2, 3, 4, 5})
util.SliceGroupBy([]string{"a", "ab", "abc"}, len) // map[1:["a"] 2:["ab"] 3:["abc"]
util.SliceCutChunks([]int{1, 2, 3, 4, 5}, 2) // [[1 2] [3 4] [5]]

// map 操作
util.MapKeys(myMap)       // []K
util.MapValues(myMap)     // []V
util.MapInvert(myMap)     // 交换键值
util.MapAssign(m1, m2, m3) // 从左到右合并

// gzip 压缩
data, _ := util.GZipCompressBytes(raw)
raw, _   := util.GZipDecompressBytes(data)
```

### writer/filewriter

**按大小滚动** — 文件大小超过限制时自动切分：

```go
w, _ := filewriter.NewSizeRollingFileWriter("./logs", "app", 10, 64<<20)
w.Write([]byte("日志行\n"))
w.Close()
```

**按时间滚动** — 按时间周期自动切分：

```go
w, _ := filewriter.NewTimeRollingFileWriter("./logs", "app", 30, filewriter.RollingPeriodDay)
w.Write([]byte("日志行\n"))
w.Close()
```

## 约定

- 全局使用 Go 泛型（1.26+）— 用 `any` 替代 `interface{}`。
- 优先使用标准库（`slices`、`maps`、`math/rand/v2`）而非手写实现。
- 函数式选项模式（`pool.Option`、cache 线程安全开关）。
- 导出符号需有文档注释。可变结构体使用指针接收者。
- 测试与源文件同目录（`_test.go`），使用标准库 `testing` + `github.com/stretchr/testify`。
- 哨兵错误：`var ErrXxx = errors.New(...)`。

## 许可证

[MIT](LICENSE)
