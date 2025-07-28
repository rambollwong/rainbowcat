package pool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewWorkerPool(t *testing.T) {
	// Test creating a worker pool with positive number of workers
	pool := NewWorkerPool(5)
	require.NotNil(t, pool)
	require.Equal(t, 5, pool.RunningWorkers())
	pool.Close()
}

func TestNewWorkerPoolWithZeroWorkers(t *testing.T) {
	// Test creating a worker pool with zero workers (should default to 1)
	pool := NewWorkerPool(0)
	require.NotNil(t, pool)
	require.Equal(t, 1, pool.RunningWorkers())
	pool.Close()
}

func TestSubmitTask(t *testing.T) {
	pool := NewWorkerPool(3)
	defer pool.Close()

	var counter int32
	task := func() {
		atomic.AddInt32(&counter, 1)
	}

	// Submit multiple tasks
	numTasks := 10
	for i := 0; i < numTasks; i++ {
		err := pool.Submit(task)
		require.NoError(t, err)
	}

	// Wait for tasks to complete
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, int32(numTasks), atomic.LoadInt32(&counter))
}

func TestSubmitTaskToClosedPool(t *testing.T) {
	pool := NewWorkerPool(1)
	pool.Close()

	task := func() {}
	err := pool.Submit(task)
	require.Error(t, err)
	require.Equal(t, "worker pool is closed", err.Error())
}

func TestSubmitTaskWithRejectHandler(t *testing.T) {
	var rejectedTasks int32
	rejectHandler := func(task Task) {
		atomic.AddInt32(&rejectedTasks, 1)
	}

	pool := NewWorkerPool(1, WithRejectHandler(rejectHandler))
	pool.Close()

	task := func() {}
	err := pool.Submit(task)
	require.Error(t, err)
	require.Equal(t, "worker pool is closed", err.Error())
	require.Equal(t, int32(1), atomic.LoadInt32(&rejectedTasks))
}

func TestCloseWithTimeout(t *testing.T) {
	pool := NewWorkerPool(2)
	defer pool.Close()

	// Submit a long-running task
	longTask := func() {
		time.Sleep(200 * time.Millisecond)
	}

	err := pool.Submit(longTask)
	require.NoError(t, err)

	// Try to close with short timeout - should fail
	result := pool.CloseWithTimeout(50 * time.Millisecond)
	require.False(t, result)

	// Try to close with longer timeout - should succeed
	result = pool.CloseWithTimeout(500 * time.Millisecond)
	require.True(t, result)
}

func TestPendingTasks(t *testing.T) {
	// Create a pool with buffered task channel
	pool := NewWorkerPool(1, WithBufferSize(5))
	defer pool.Close()

	// Submit tasks that take some time to complete
	blockingTask := func() {
		time.Sleep(100 * time.Millisecond)
	}

	// Submit more tasks than workers
	for i := 0; i < 3; i++ {
		err := pool.Submit(blockingTask)
		require.NoError(t, err)
	}

	// Give some time for tasks to be picked up by workers
	time.Sleep(10 * time.Millisecond)

	// Check pending tasks count
	pending := pool.PendingTasks()
	require.GreaterOrEqual(t, pending, 0)
	require.Less(t, pending, 3)
}

func TestRunningWorkers(t *testing.T) {
	pool := NewWorkerPool(5)
	require.Equal(t, 5, pool.RunningWorkers())

	pool.Close()
	require.Equal(t, 0, pool.RunningWorkers())
}

func TestWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a pool with a context that will be cancelled
	pool := NewWorkerPool(2, WithContext(ctx))
	defer pool.Close()

	var counter int32
	task := func() {
		atomic.AddInt32(&counter, 1)
	}

	// Submit a task
	err := pool.Submit(task)
	require.NoError(t, err)

	// Cancel the context
	cancel()

	// Submit another task after cancellation
	err = pool.Submit(task)
	require.Error(t, err)
	require.Equal(t, "worker pool is closing", err.Error())

	// Wait and check that only the first task was executed
	time.Sleep(50 * time.Millisecond)
	require.Equal(t, int32(1), atomic.LoadInt32(&counter))
}

func TestWorkerPool_RunningWorkers(t *testing.T) {
	pool := NewWorkerPool(5)
	defer pool.Close()

	// Should report the correct number of running workers
	running := pool.RunningWorkers()
	require.Equal(t, 5, running)

	// Close the pool
	pool.Close()

	// Should report 0 running workers after closing
	running = pool.RunningWorkers()
	require.Equal(t, 0, running)
}
