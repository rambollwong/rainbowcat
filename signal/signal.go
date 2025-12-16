package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// OsSignalsC is a channel for receiving OS signals
var (
	OsSignalsC chan os.Signal
)

// init initializes the OsSignalsC channel with a buffer size of 2
func init() {
	OsSignalsC = make(chan os.Signal, 2)
}

// SendExitSignal sends an interrupt signal to the OsSignalsC channel
func SendExitSignal() {
	OsSignalsC <- os.Interrupt
}

// WatchExitSignal watches for SIGINT and SIGTERM signals and executes the callback function when received
func WatchExitSignal(callback func()) {
	go func() {
		// Notify OsSignalsC channel of incoming SIGINT and SIGTERM signals
		signal.Notify(OsSignalsC, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-OsSignalsC:
			// Execute the callback function when a signal is received
			callback()
		}
	}()
}

// WatchExitSignalWithContext watches for SIGINT and SIGTERM signals and executes the callback function when received.
// It also respects the context cancellation
// - if context is cancelled, it will stop watching and return immediately.
func WatchExitSignalWithContext(ctx context.Context, callback func()) {
	go func() {
		// Notify OsSignalsC channel of incoming SIGINT and SIGTERM signals
		signal.Notify(OsSignalsC, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			// Context was cancelled, exit the goroutine
			return
		case <-OsSignalsC:
			// Execute the callback function when a signal is received
			callback()
		}
	}()
}
