package time

import "time"

// UseLocalUTC sets the global time.Local to time.UTC.
//
// WARNING: This mutates the process-wide time.Local variable. All code that
// calls time.Now(), time.Parse(), or any other function that depends on the
// local timezone — including third-party libraries — will be affected.
// Only call this at program startup when you are certain no other code
// depends on the system-local timezone.
func UseLocalUTC() {
	time.Local = time.UTC
}
