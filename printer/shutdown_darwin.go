//go:build darwin

package printer

// Shutdown releases printer resources on agent stop.
// Bluetooth connections are opened and closed per print job; nothing to tear down here.
func Shutdown() {}
