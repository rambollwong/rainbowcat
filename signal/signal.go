package signal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"golang.org/x/sys/unix"
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

// WaitForKeyPress waits for the user to press any key to continue
func WaitForKeyPress() {
	fmt.Println("\nPress any key to exit...")

	if runtime.GOOS == "windows" {
		// Windows processing (as above, using bufio or Windows API)
		_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	} else {
		// Linux/macOS processing (using x/sys/unix package)
		// Save original terminal settings
		fd := int(os.Stdin.Fd())
		oldState, err := unix.IoctlGetTermios(fd, unix.TCGETS) // Replace syscall.Gettermios
		if err != nil {
			fmt.Println("Failed to get terminal settings, exiting automatically...")
			return
		}
		newState := *oldState // Copy original settings

		// Modify terminal settings: disable canonical mode (no need for Enter), disable echo
		newState.Lflag &^= unix.ICANON // Disable canonical mode
		newState.Lflag &^= unix.ECHO   // Disable echo
		newState.Cc[unix.VMIN] = 1     // Read at least 1 character
		newState.Cc[unix.VTIME] = 0    // No timeout

		// Apply new settings
		err = unix.IoctlSetTermios(fd, unix.TCSETS, &newState) // Replace syscall.Settermios
		if err != nil {
			fmt.Println("Failed to modify terminal settings, exiting automatically...")
			return
		}
		defer unix.IoctlSetTermios(fd, unix.TCSETS, oldState) // Restore original settings

		// Read one character (any key)
		var buf [1]byte
		_, _ = os.Stdin.Read(buf[:])
	}
}
