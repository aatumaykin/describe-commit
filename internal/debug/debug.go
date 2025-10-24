package debug

import (
	"os"
	"strconv"
	"sync/atomic"
)

var Enabled atomic.Bool //nolint:gochecknoglobals

// SetColorEnabled sets the color enabled state for debug output
func SetColorEnabled(enabled bool) {
	SetColorEnabledState(enabled)
}

// Printf is a helper function to print debug information to the stderr.
func Printf(format string, args ...any) {
	if Enabled.Load() {
		DebugColorPrintf(format, args...)
	}
}

func init() { //nolint:gochecknoinits
	const debugEnvName = "DEBUG" // environment variable name to enable debug output

	if v, ok := os.LookupEnv(debugEnvName); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			Enabled.Store(b)
		}
	}
}
