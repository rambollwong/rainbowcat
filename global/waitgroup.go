package global

import (
	"context"
	"sync"
)

var (
	wg sync.WaitGroup // Global wait group, used to control the program to exit gracefully
)

func Wait() {
	wg.Wait()
}

func RunTask(task func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		task()
	}()
}
func RunTaskWithContext(ctx context.Context, task func(ctx context.Context)) {
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		task(ctx)
	}(ctx)
}
