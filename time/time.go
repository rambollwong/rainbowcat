package time

import "time"

func UseLocalUTC() {
	// Use UTC time globally
	time.Local = time.UTC
}
