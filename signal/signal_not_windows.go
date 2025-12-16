//go:build linux || darwin || freebsd

package signal

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// WaitForKeyPress waits for the user to press any key to continue
func WaitForKeyPress() {
	fmt.Println("\nPress any key to exit...")

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
