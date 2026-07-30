# 🐱 rainbowcat

A useful library for Go.

[🇨🇳 中文版本](README_zh_cn.md)

## Installation

```bash
go get github.com/rambollwong/rainbowcat
```

Requires **Go 1.26+**.

## Packages

| Package | Description |
|---|---|
| [`cache`](#cache) | Generic FIFO cache with optional thread-safety and eviction callback |
| [`global`](#global) | Global `sync.WaitGroup` for graceful shutdown coordination |
| [`logo`](#logo) | ASCII art banner printer (cosmetic) |
| [`pipeline`](#pipeline) | Multi-stage parallel task pipeline — each stage runs concurrently |
| [`pool`](#pool) | Goroutine worker pool + `[]byte` reuse pool |
| [`signal`](#signal) | OS signal handling (`SIGINT`/`SIGTERM`) with platform variants |
| [`smtp`](#smtp) | SMTP email sender (wraps `gomail`) |
| [`task`](#task) | Timer/ticker-based task monitor with pluggable `DataStore` |
| [`time`](#time) | Single helper: `UseLocalUTC()` |
| [`types`](#types) | Generic helpers: `Set[T]`, `Entry[K,V]`, `Clonable[T]` |
| [`util`](#util) | Slice & map operations, bytes/hex, gzip, fibonacci, string parsing |
| [`writer/filewriter`](#writerfilewriter) | Rolling file writers: size-based and time-based |

---

### cache

A generic fixed-size FIFO cache backed by `container/list` + `map`, with optional `sync.RWMutex` thread-safety and eviction callback.

```go
c := cache.NewFIFOCache[int, string](100, true)
c.Put(1, "hello")
v, ok := c.Get(1) // "hello", true
c.Size()           // 1
```

### global

Process-global `sync.WaitGroup` for tracking running goroutines during graceful shutdown.

```go
global.RunTask(func() {
    // do work
})
// ... wait for all tasks to finish
global.Wait()
```

### pool

**Worker pool** — bounded goroutine pool with functional options:

```go
wp := pool.NewWorkerPool(4,
    pool.WithBufferSize(16),
    pool.WithRejectHandler(func(t pool.Task) { /* log */ }),
)
err := wp.Submit(func() { /* do work */ })
wp.Close()
```

**Bytes pool** — `sync.Pool` wrapper for `[]byte` reuse:

```go
bp := pool.NewBytesPool(512, pool.DefaultMaxBytesCap)
bz := bp.Get()
defer bp.Put(bz)
```

### pipeline

Multi-stage concurrent pipeline. Each stage has configurable concurrency; output of stage N becomes input of stage N+1.

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

Watch for `SIGINT`/`SIGTERM` and execute callbacks:

```go
signal.WatchExitSignal(func() {
    fmt.Println("shutting down...")
})
signal.WaitForKeyPress() // cross-platform "press any key"
```

### task

Timer/ticker scheduling with a pluggable `DataStore`:

```go
tm := &task.TasksMonitor{}
tm.SetDataStore(myDataStore)
tm.RegisterTickerForTasks(time.Minute, "cleanup", func(d task.Data) {
    // periodic cleanup
})
tm.Start()
```

### types

Thread-safe generic `Set[T]` backed by `map + sync.RWMutex`:

```go
s := types.NewSet[int]()
s.Put(1)      // true
s.Put(1)      // false (already exists)
s.Exist(1)    // true
s.Size()      // 1
s.Range(func(v int) bool { fmt.Println(v); return true })
```

### util

Rich collection of slice and map utilities (inspired by `samber/lo`), plus gzip compression, fibonacci, byte conversions, and string parsing:

```go
// slice operations
util.SliceFilter([]int{1, 2, 3, 4}, func(i int, v int) bool { return v%2 == 0 })
util.SliceShuffle([]int{1, 2, 3, 4, 5})
util.SliceGroupBy([]string{"a", "ab", "abc"}, len) // map[1:["a"] 2:["ab"] 3:["abc"]
util.SliceCutChunks([]int{1, 2, 3, 4, 5}, 2) // [[1 2] [3 4] [5]]

// map operations
util.MapKeys(myMap)       // []K
util.MapValues(myMap)     // []V
util.MapInvert(myMap)     // swap keys ↔ values
util.MapAssign(m1, m2, m3) // merge left-to-right

// gzip
data, _ := util.GZipCompressBytes(raw)
raw, _   := util.GZipDecompressBytes(data)
```

### writer/filewriter

**Size-based rolling** — rotates when file exceeds a size limit:

```go
w, _ := filewriter.NewSizeRollingFileWriter("./logs", "app", 10, 64<<20)
w.Write([]byte("log line\n"))
w.Close()
```

**Time-based rolling** — rotates at time-period boundaries:

```go
w, _ := filewriter.NewTimeRollingFileWriter("./logs", "app", 30, filewriter.RollingPeriodDay)
w.Write([]byte("log line\n"))
w.Close()
```

## Conventions

- Go generics everywhere (1.26+) — `any` over `interface{}`.
- Prefer standard library (`slices`, `maps`, `math/rand/v2`) over hand-rolled equivalents.
- Functional options pattern (`pool.Option`, cache thread-safety toggle).
- Exported symbols have doc comments. Pointer receivers for mutable structs.
- Tests colocated (`_test.go`), using standard `testing` + `github.com/stretchr/testify`.
- Sentinel errors: `var ErrXxx = errors.New(...)`.

## License

[MIT](LICENSE)
