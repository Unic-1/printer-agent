//go:build !darwin

package printer

// Shutdown releases printer resources on agent stop. No-op on non-macOS platforms.
func Shutdown() {}
