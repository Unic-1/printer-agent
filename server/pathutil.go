package server

import (
	"os"
	"path/filepath"
)

// ExeDir returns the directory containing the currently running executable.
// This is critical for Windows Service mode where the working directory
// defaults to C:\Windows\System32, making relative paths useless.
func ExeDir() string {
	exe, err := os.Executable()
	if err != nil {
		// Fallback to current working directory if os.Executable fails
		return "."
	}
	return filepath.Dir(exe)
}
