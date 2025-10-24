package debug

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
)

// ColorEnabled is a global flag that controls whether colors should be used
var ColorEnabled bool

// init checks if colors should be enabled based on environment and terminal capabilities
func init() {
	ColorEnabled = shouldEnableColors()
}

// shouldEnableColors determines if colors should be enabled based on:
// 1. NO_COLOR environment variable
// 2. TERM environment variable
// 3. Platform-specific checks
func shouldEnableColors() bool {
	// Check for NO_COLOR environment variable (https://no-color.org/)
	if _, exists := os.LookupEnv("NO_COLOR"); exists {
		return false
	}

	// Check if we're running on Windows without ANSI support
	if runtime.GOOS == "windows" {
		// Check for Windows Terminal or other modern terminals
		term := os.Getenv("TERM")
		if term == "" {
			// Check for Windows Terminal
			wtSession := os.Getenv("WT_SESSION")
			if wtSession == "" {
				// Check for ConEmu
				conEmuANSI := os.Getenv("ConEmuANSI")
				if conEmuANSI != "ON" {
					return false
				}
			}
		}
	}

	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "" {
		return false
	}

	// Check if terminal supports colors
	if !strings.Contains(term, "color") &&
		!strings.Contains(term, "xterm") &&
		!strings.Contains(term, "screen") &&
		!strings.Contains(term, "tmux") {
		return false
	}

	// Check if output is redirected
	if !isTerminal(os.Stderr) {
		return false
	}

	return true
}

// isTerminal checks if the given file descriptor is a terminal
func isTerminal(f *os.File) bool {
	// Simple check: if we can get file info and it's a character device
	stat, err := f.Stat()
	if err != nil {
		return false
	}

	// Check if it's a character device (terminal)
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// Colorize applies color to text if colors are enabled
func Colorize(color, text string) string {
	if !ColorEnabled {
		return text
	}
	return color + text + ColorReset
}

// ColorPrintf prints formatted text with color if colors are enabled
func ColorPrintf(color, format string, args ...any) {
	if !ColorEnabled {
		fmt.Fprintf(os.Stderr, format, args...)
		return
	}

	formatted := fmt.Sprintf(format, args...)
	colored := Colorize(color, formatted)
	fmt.Fprint(os.Stderr, colored)
}

// DebugColorPrintf prints debug information with color
func DebugColorPrintf(format string, args ...any) {
	if !Enabled.Load() {
		return
	}

	debugMsg := fmt.Sprintf("# [debug] %s\n", fmt.Sprintf(format, args...))
	ColorPrintf(ColorGray, "%s", debugMsg)
}

// DebugHeaderColorPrintf prints debug headers with a different color
func DebugHeaderColorPrintf(format string, args ...any) {
	if !Enabled.Load() {
		return
	}

	debugMsg := fmt.Sprintf("# [debug] %s\n", fmt.Sprintf(format, args...))
	ColorPrintf(ColorCyan, "%s", debugMsg)
}

// DebugContentColorPrintf prints debug content with a muted color
func DebugContentColorPrintf(format string, args ...any) {
	if !Enabled.Load() {
		return
	}

	debugMsg := fmt.Sprintf("%s\n", fmt.Sprintf(format, args...))
	ColorPrintf(ColorGray, "%s", debugMsg)
}

// InfoColorPrintf prints info messages with color
func InfoColorPrintf(format string, args ...any) {
	ColorPrintf(ColorCyan, format, args...)
}

// DryRunHeaderColorPrintf prints dry run headers with color
func DryRunHeaderColorPrintf(format string, args ...any) {
	ColorPrintf(ColorYellow, format, args...)
}

// DryRunContentColorPrintf prints dry run content with color
func DryRunContentColorPrintf(format string, args ...any) {
	ColorPrintf(ColorWhite, format, args...)
}

// SuccessColorPrintf prints success messages with color
func SuccessColorPrintf(format string, args ...any) {
	ColorPrintf(ColorGreen, format, args...)
}

// WarningColorPrintf prints warning messages with color
func WarningColorPrintf(format string, args ...any) {
	ColorPrintf(ColorYellow, format, args...)
}

// ErrorColorPrintf prints error messages with color
func ErrorColorPrintf(format string, args ...any) {
	ColorPrintf(ColorRed, format, args...)
}

// SetColorEnabledState manually sets the color enabled state
func SetColorEnabledState(enabled bool) {
	ColorEnabled = enabled
}

// IsColorEnabled returns whether colors are currently enabled
func IsColorEnabled() bool {
	return ColorEnabled
}
