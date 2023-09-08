package pool

import (
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
)

func TestNewBytesPool(t *testing.T) {
	NewBytesPool(DefaultNewBytesCap, DefaultMaxBytesCap)
}

func TestBytesPoolGet(t *testing.T) {
	p := NewBytesPool(DefaultNewBytesCap, DefaultMaxBytesCap)
	bz := p.Get()
	require.Equal(t, 0, len(*bz))
	require.Equal(t, DefaultNewBytesCap, cap(*bz))
}

func TestBytesPoolPut(t *testing.T) {
	var newTimes uint64
	p := NewBytesPool(3, DefaultMaxBytesCap)
	tmp := p.p.New
	p.p.New = func() interface{} {
		atomic.AddUint64(&newTimes, 1)
		return tmp()
	}

	bz := p.Get()
	*bz = append(*bz, byte(1), byte(2), byte(3))
	p.Put(bz)

	bz2 := p.Get()
	require.True(t, newTimes == 1)
	require.True(t, cap(*bz2) == 3)
	require.True(t, len(*bz2) == 0)

	*bz2 = append(*bz2, byte(1), byte(2), byte(3), byte(4))
	p.Put(bz2)
	bz3 := p.Get()
	require.True(t, newTimes == 1)
	require.True(t, cap(*bz3) > 3)
	require.True(t, len(*bz3) == 0)

}
