//go:build darwin

package printer

import (
	"os"
	"os/exec"
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
	var printers []*models.Printer

	// 1. Serial device discovery (USB-to-serial adapters)
	devPrinters, _ := discoverSerialDevices()
	printers = append(printers, devPrinters...)

	// 2. CUPS-registered USB printers
	cupsPrinters, _ := discoverCUPSUSBPrinters()
	printers = append(printers, cupsPrinters...)

	return printers, nil
}

func discoverSerialDevices() ([]*models.Printer, error) {
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

// discoverCUPSUSBPrinters finds printers connected via USB that are registered
// with CUPS. Uses `lpstat -v` to get device URIs and filters for usb:// ones.
func discoverCUPSUSBPrinters() ([]*models.Printer, error) {
	// lpstat -v outputs lines like:
	//   device for PrinterName: usb://EPSON/TM-T88V?serial=...
	//   device for PDF_Printer: cups-pdf:/
	out, err := exec.Command("lpstat", "-v").Output()
	if err != nil {
		return nil, err
	}

	var printers []*models.Printer
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "device for ") {
			continue
		}

		// Parse "device for <name>: <uri>"
		rest := strings.TrimPrefix(line, "device for ")
		colonIdx := strings.Index(rest, ": ")
		if colonIdx < 0 {
			continue
		}
		name := rest[:colonIdx]
		uri := rest[colonIdx+2:]

		if !strings.HasPrefix(uri, "usb://") {
			continue
		}

		printers = append(printers, &models.Printer{
			ID:      "usb:printer:" + name,
			Name:    name + " (USB)",
			Type:    models.PrinterUSB,
			Address: "printer:" + name,
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
