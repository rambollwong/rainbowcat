package pool

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Task represents a unit of work to be executed
type Task func()

// WorkerPool manages a pool of goroutines to execute tasks concurrently
type WorkerPool struct {
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	tasks         chan Task
	workers       int
	running       bool
	mu            sync.RWMutex
	rejectHandler func(Task) // Handler for rejected tasks
}

// Option is a functional option for configuring the worker pool
type Option func(pool *WorkerPool)

// WithRejectHandler sets the handler function for rejected tasks
func WithRejectHandler(handler func(Task)) Option {
	return func(pool *WorkerPool) {
		pool.rejectHandler = handler
	}
}

// WithBufferSize sets the buffer size for the task channel
func WithBufferSize(size int) Option {
	return func(pool *WorkerPool) {
		if size > 0 {
			pool.tasks = make(chan Task, size)
		}
	}
}

// WithContext sets the context for the worker pool.
// The provided context will be used to control the lifecycle of the worker pool.
// When the context is cancelled, all workers will begin shutting down.
//
//	ctx : the context to use for controlling worker pool lifecycle
func WithContext(ctx context.Context) Option {
	return func(pool *WorkerPool) {
		pool.ctx, pool.cancel = context.WithCancel(ctx)
	}
}

// New creates a new worker pool with the specified number of workers
func New(workers int, opts ...Option) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		ctx:     ctx,
		cancel:  cancel,
		workers: workers,
		tasks:   make(chan Task), // Unbuffered channel by default
		running: true,
	}

	// Apply configuration options
	for _, opt := range opts {
		opt(pool)
	}

	// Start worker goroutines
	pool.startWorkers()
	return pool
}

// startWorkers initializes and starts the worker goroutines
func (p *WorkerPool) startWorkers() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-p.ctx.Done():
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					task()
				}
			}
		}()
	}
}

// Submit adds a task to the worker pool for execution
func (p *WorkerPool) Submit(task Task) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.running {
		if p.rejectHandler != nil {
			p.rejectHandler(task)
		}
		return errors.New("worker pool is closed")
	}

	select {
	case <-p.ctx.Done():
		if p.rejectHandler != nil {
			p.rejectHandler(task)
		}
		return errors.New("worker pool is closing")
	case p.tasks <- task:
		return nil
	}
}

// Close gracefully shuts down the worker pool
func (p *WorkerPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return
	}

	p.running = false
	p.cancel() // Signal all workers to exit
	close(p.tasks)
	p.wg.Wait() // Wait for all workers to complete
}

// CloseWithTimeout shuts down the worker pool with a timeout
func (p *WorkerPool) CloseWithTimeout(timeout time.Duration) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return true
	}

	p.running = false
	p.cancel() // Signal all workers to exit

	done := make(chan struct{})
	go func() {
		close(p.tasks)
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// RunningWorkers returns the number of active workers
func (p *WorkerPool) RunningWorkers() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.running {
		return 0
	}

	// Note: This is an approximation, actual count may vary during shutdown
	return p.workers
}

// PendingTasks returns the number of tasks waiting to be processed
func (p *WorkerPool) PendingTasks() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.running {
		return 0
	}

	return len(p.tasks)
}
