package pool

import "sync"

const (
	// DefaultMaxBytesCap is the default max cap value of the bytes slice in pool,
	// when Put method is called, if the cap of the bytes slice put is greater than max value,
	// it will be dropped.
	DefaultMaxBytesCap = 16 << 10 // 16k
	// DefaultNewBytesCap is the default cap value of the new bytes slice created by pool.
	DefaultNewBytesCap = 512
)

var globalBytesPool *BytesPool
var bytesPoolOnce sync.Once

// BytesPool is a pool provides bytes slice.
type BytesPool struct {
	p       *sync.Pool
	maxCap  int
	initCap int
}

// NewBytesPool create a new BytesPool instance.
//
//	initCap : the cap value of the new bytes slice created by pool.
//	maxCap  : the max cap value of the bytes slice in pool.
func NewBytesPool(initCap, maxCap int) *BytesPool {
	if initCap < 1 {
		initCap = DefaultNewBytesCap
	}
	if maxCap < initCap {
		maxCap = initCap
	}
	return &BytesPool{
		p: &sync.Pool{
			New: func() interface{} {
				bz := make([]byte, 0, initCap)
				return &bz
			},
		},
		maxCap:  maxCap,
		initCap: initCap,
	}
}

// Get borrows a bytes slice from pool. If the pool is empty, the new bytes slice will be created and returned.
func (p *BytesPool) Get() *[]byte {
	return p.p.Get().(*[]byte)
}

// Put take a bytes slice back to the pool. If the cap of the bytes slice is greater than max value, drop it.
func (p *BytesPool) Put(bz *[]byte) {
	if cap(*bz) > p.maxCap {
		return
	}
	b := (*bz)[:0]
	p.p.Put(&b)
}

func initGlobalBytesPool() {
	bytesPoolOnce.Do(func() {
		globalBytesPool = NewBytesPool(DefaultNewBytesCap, DefaultMaxBytesCap)
	})
}

// BytesPoolPut take a bytes slice back to the global pool.
// If the cap of the bytes slice is greater than max value, drop it.
// The max cap value default DefaultMaxBytesCap.
func BytesPoolPut(bz *[]byte) {
	initGlobalBytesPool()
	globalBytesPool.Put(bz)
}

// BytesPoolGet borrows a bytes slice from global pool.
// If the pool is empty, the new bytes slice will be created and returned.
func BytesPoolGet() *[]byte {
	initGlobalBytesPool()
	return globalBytesPool.Get()
}

// SetBytesPoolMaxCap set the max cap for global pool.
func SetBytesPoolMaxCap(maxCap int) {
	globalBytesPool.maxCap = maxCap
}
