//go:build windows

package signal

import (
	"bufio"
	"fmt"
	"os"
)

// WaitForKeyPress waits for the user to press any key to continue
func WaitForKeyPress() {
	fmt.Println("\nPress enter to exit...")

	// Windows processing (as above, using bufio or Windows API)
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')

}
