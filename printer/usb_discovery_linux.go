//go:build linux

package printer

import (
	"os"
	"path/filepath"
	"strings"

	"printer-agent/models"
)

func DiscoverUSB() ([]*models.Printer, error) {
	var printers []*models.Printer

	// CUPS / kernel USB printer nodes
	for _, pattern := range []string{"/dev/usb/lp*", "/dev/lp*"} {
		matches, _ := filepath.Glob(pattern)
		for _, path := range matches {
			base := filepath.Base(path)
			printers = append(printers, &models.Printer{
				ID:      "usb:" + path,
				Name:    base,
				Type:    models.PrinterUSB,
				Address: path,
				Online:  true,
			})
		}
	}

	// USB serial adapters (same pattern as macOS)
	entries, err := os.ReadDir("/dev")
	if err == nil {
		for _, e := range entries {
			name := e.Name()
			if strings.HasPrefix(name, "ttyUSB") || strings.HasPrefix(name, "ttyACM") {
				path := "/dev/" + name
				printers = append(printers, &models.Printer{
					ID:      "usb:" + path,
					Name:    name,
					Type:    models.PrinterUSB,
					Address: path,
					Online:  true,
				})
			}
		}
	}

	return printers, nil
}
