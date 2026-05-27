package printer

import "strings"

// NormalizeUSBAddress converts discovery IDs and legacy saved values into a
// value suitable for printUSB (COM3, USB001, or printer:Name on Windows).
func NormalizeUSBAddress(addr string) string {
	a := strings.TrimSpace(addr)
	if a == "" {
		return a
	}

	if strings.HasPrefix(a, "usb:printer:") {
		return "printer:" + strings.TrimPrefix(a, "usb:printer:")
	}
	if strings.HasPrefix(a, "usb:port:") {
		return strings.TrimPrefix(a, "usb:port:")
	}
	if strings.HasPrefix(a, "usb:") {
		rest := strings.TrimPrefix(a, "usb:")
		if strings.HasPrefix(rest, "printer:") || strings.HasPrefix(rest, "port:") {
			return NormalizeUSBAddress(rest)
		}
		return rest
	}

	return a
}
