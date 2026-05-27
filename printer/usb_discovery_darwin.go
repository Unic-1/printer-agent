//go:build darwin

package printer

import (
	"os"
	"strings"

	"printer-agent/models"
)

// USB serial / direct paths commonly used by ESC/POS printers on macOS.
var darwinUSBDevicePrefixes = []string{
	"cu.usbserial",
	"cu.usbmodem",
	"cu.usb",
	"cu.SLAB_USB",
	"cu.wchusbserial",
}

func DiscoverUSB() ([]*models.Printer, error) {
	entries, err := os.ReadDir("/dev")
	if err != nil {
		return nil, err
	}

	var printers []*models.Printer
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "cu.") {
			continue
		}
		if !isDarwinUSBDevice(name) {
			continue
		}
		path := "/dev/" + name
		printers = append(printers, &models.Printer{
			ID:      "usb:" + path,
			Name:    name,
			Type:    models.PrinterUSB,
			Address: path,
			Online:  true,
		})
	}
	return printers, nil
}

func isDarwinUSBDevice(name string) bool {
	for _, prefix := range darwinUSBDevicePrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
